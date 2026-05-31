package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

var (
	ErrFileTooLarge        = errors.New("文件大小超出限制 (最大 200MB)")
	ErrFileTypeUnsupported = errors.New("不支持的文件格式，仅支持 txt/epub/mobi/pdf")
	ErrFileEmpty           = errors.New("文件为空")
	ErrUploadNotFound      = errors.New("上传记录不存在")
	ErrChapterNotFound     = errors.New("章节不存在")
	ErrRuleNotFound        = errors.New("规则不存在")
	ErrFilterRuleNotFound  = errors.New("过滤规则不存在")
	ErrBookFileNotFound    = errors.New("文件记录不存在")
	ErrFileMD5Exists       = errors.New("文件 MD5 已存在（重复文件）")
)

const (
	MaxFileSize       = 200 * 1024 * 1024 // 200MB
	StorageBaseDir    = "storage/books"
	DefaultChapterLen = 100
	MaxChapterLen     = 100000
)

// BookFileService 小说文件管理服务
type BookFileService struct {
	db              *gorm.DB
	bookRepo        *repository.BookRepository
	bookFileRepo    *repository.BookFileRepository
	uploadRepo      *repository.BookUploadRepository
	chapterRepo     *repository.BookChapterRepository
	chapterRuleRepo *repository.BookChapterRuleRepository
	filterRuleRepo  *repository.BookContentFilterRuleRepository
	categoryRepo    *repository.BookCategoryRepository
	tagRepo         *repository.BookTagRepository
}

func NewBookFileService(
	db *gorm.DB,
	bookRepo *repository.BookRepository,
	bookFileRepo *repository.BookFileRepository,
	uploadRepo *repository.BookUploadRepository,
	chapterRepo *repository.BookChapterRepository,
	chapterRuleRepo *repository.BookChapterRuleRepository,
	filterRuleRepo *repository.BookContentFilterRuleRepository,
	categoryRepo *repository.BookCategoryRepository,
	tagRepo *repository.BookTagRepository,
) *BookFileService {
	return &BookFileService{
		db:              db,
		bookRepo:        bookRepo,
		bookFileRepo:    bookFileRepo,
		uploadRepo:      uploadRepo,
		chapterRepo:     chapterRepo,
		chapterRuleRepo: chapterRuleRepo,
		filterRuleRepo:  filterRuleRepo,
		categoryRepo:    categoryRepo,
		tagRepo:         tagRepo,
	}
}

// ==================== 文件上传 ====================

// Upload 上传文件并创建上传记录，返回建议的书名和作者
func (s *BookFileService) Upload(ctx context.Context, reader io.Reader, originalName string, fileSize uint64, userID uint64) (*dto.FileUploadResponse, error) {
	// 1. 格式验证
	ext := strings.ToLower(filepath.Ext(originalName))
	if !ValidateFileType(ext) {
		return nil, ErrFileTypeUnsupported
	}
	// 2. 大小验证
	if fileSize > MaxFileSize {
		return nil, ErrFileTooLarge
	}
	// 3. 读取内容并计算 MD5
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	if len(data) == 0 {
		return nil, ErrFileEmpty
	}
	md5Hash := fmt.Sprintf("%x", md5.Sum(data))

	// 4. 写入本地存储
	relPath := filepath.Join(StorageBaseDir, fmt.Sprintf("%s_%s", md5Hash[:8], originalName))
	absPath := filepath.Join(".", relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}
	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	// 5. 从文件内容提取书名和作者
	title, author := extractMetaFromContent(data, originalName)

	// 6. 匹配数据库中已存在的小说（同名同作者）
	var matchedID *uint64
	var matchedTitle string
	if title != "" {
		book, findErr := s.bookRepo.FindByTitleAndAuthor(ctx, title, author)
		if findErr == nil && book != nil {
			matchedID = &book.ID
			matchedTitle = book.Title
		}
	}

	// 7. 创建上传记录
	sourceFormat := strings.TrimPrefix(ext, ".")
	upload := &model.BookUpload{
		OriginalName: originalName,
		FilePath:     relPath,
		FileSize:     fileSize,
		FileMD5:      &md5Hash,
		SourceFormat: &sourceFormat,
		ParseStatus:  model.ParsePending,
	}
	upload.CreateBy = &userID
	upload.UpdateBy = &userID
	if err := s.uploadRepo.Create(ctx, upload); err != nil {
		return nil, fmt.Errorf("创建上传记录失败: %w", err)
	}

	return &dto.FileUploadResponse{
		UploadID:         upload.ID,
		OriginalName:     originalName,
		FileSize:         fileSize,
		SourceFormat:     &sourceFormat,
		SuggestedTitle:   title,
		SuggestedAuthor:  author,
		MatchedBookID:    matchedID,
		MatchedBookTitle: matchedTitle,
	}, nil
}

// ConfirmImport 用户确认入库：匹配或创建 Book，写入 BookFile + 章节索引
func (s *BookFileService) ConfirmImport(ctx context.Context, req *dto.ConfirmImportRequest, userID uint64) (*dto.ConfirmImportResponse, error) {
	up, err := s.uploadRepo.GetByID(ctx, req.UploadID)
	if err != nil {
		return nil, ErrUploadNotFound
	}
	if up.ParseStatus != model.ParsePending {
		return nil, errors.New("该上传记录已处理，不可重复入库")
	}

	// 读取文件
	data, err := os.ReadFile(up.FilePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 获取章节识别规则（全局 + 书籍级）
	rules, _ := s.chapterRuleRepo.ListEffective(ctx, 0)

	// 先应用入库过滤规则，再在过滤后的数据上解析章节
	// 确保字节偏移与存储到文件的内容一致
	filterRules, _ := s.filterRuleRepo.ListByStage(ctx, model.FilterStageInput)
	filter := NewContentFilter(filterRules)
	filteredData := applyFilterToContent(data, filter)

	parser := NewChapterParser(rules)
	parseResult := parser.Parse(filteredData)

	// 事务入库
	var bookID uint64
	var chapterCount uint32
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 匹配已有书籍或新建
		book, findErr := s.bookRepo.FindByTitleAndAuthor(ctx, req.Title, req.Author)
		if findErr != nil {
			book = &model.Book{
				Title:         req.Title,
				Author:        req.Author,
				SerialStatus:  model.SerialOngoing,
				Visibility:    model.VisibilityPublic,
				TotalChapters: uint32(len(parseResult.Chapters)),
				TotalWords:    sumWords(parseResult.Chapters),
				Status:        model.BookReviewing,
			}
			book.CreateBy = &userID
			book.UpdateBy = &userID
			if err := tx.Create(book).Error; err != nil {
				return fmt.Errorf("创建书籍失败: %w", err)
			}
		}
		bookID = book.ID

		// 写入解析后的内容文件
		contentRelPath := filepath.Join(StorageBaseDir, fmt.Sprintf("content_%d.txt", bookID))
		contentAbsPath := filepath.Join(".", contentRelPath)
		if err := os.MkdirAll(filepath.Dir(contentAbsPath), 0755); err != nil {
			return fmt.Errorf("创建内容目录失败: %w", err)
		}
		if err := os.WriteFile(contentAbsPath, filteredData, 0644); err != nil {
			return fmt.Errorf("写入内容文件失败: %w", err)
		}

		// 创建 BookFile 记录
		contentMD5 := fmt.Sprintf("%x", md5.Sum(filteredData))
		charset := DetectCharset(filteredData)
		bf := &model.BookFile{
			BookID:         bookID,
			OriginalName:   up.OriginalName,
			SourceType:     model.FileSourceUserUpload,
			SourceFormat:   up.SourceFormat,
			SourceFileURL:  &up.FilePath,
			ContentPath:    &contentRelPath,
			ContentSize:    uint64(len(filteredData)),
			ContentMD5:     &contentMD5,
			ContentCharset: charset,
			ContentVersion: 1,
			ChapterCount:   uint32(len(parseResult.Chapters)),
			IsPrimary:      true,
			FileStatus:     model.FileSuccess,
		}
		bf.CreateBy = &userID
		bf.UpdateBy = &userID
		if err := tx.Create(bf).Error; err != nil {
			return fmt.Errorf("创建文件记录失败: %w", err)
		}

		// 创建章节索引
		chapters := make([]model.BookChapter, len(parseResult.Chapters))
		for i, seg := range parseResult.Chapters {
			chapters[i] = model.BookChapter{
				BookID:     bookID,
				FileID:     bf.ID,
				ChapterNo:  uint32(i + 1),
				Title:      seg.Title,
				ByteOffset: seg.ByteOffset,
				ByteLength: seg.ByteLength,
				WordCount:  seg.WordCount,
				Status:     model.ChapterPublished,
			}
		}
		if err := tx.Create(&chapters).Error; err != nil {
			return fmt.Errorf("创建章节索引失败: %w", err)
		}

		// 更新书籍的主文件ID和统计
		book.PrimaryFileID = &bf.ID
		book.TotalChapters = uint32(len(parseResult.Chapters))
		book.TotalWords = sumWords(parseResult.Chapters)
		book.UpdateBy = &userID
		if err := tx.Save(book).Error; err != nil {
			return err
		}

		// 更新上传记录
		chapterCount = uint32(len(parseResult.Chapters))
		up.BookID = &bookID
		up.ParseStatus = model.ParseSuccess
		up.ChapterCount = &chapterCount
		up.UpdateBy = &userID
		if err := tx.Save(up).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.ConfirmImportResponse{
		UploadID:     up.ID,
		BookID:       bookID,
		BookTitle:    req.Title,
		BookAuthor:   req.Author,
		ChapterCount: chapterCount,
		ParseStatus:  string(model.ParseSuccess),
	}, nil
}

// ==================== 扫描入库 ====================

// ScanPending 扫描所有待处理的上传任务
func (s *BookFileService) ScanPending(ctx context.Context) (*dto.ScanAllResponse, error) {
	pending, err := s.uploadRepo.ListByParseStatus(ctx, model.ParsePending, 50)
	if err != nil {
		return nil, err
	}
	resp := &dto.ScanAllResponse{Results: make([]dto.ScanResult, 0, len(pending))}
	for _, up := range pending {
		result := s.scanSingle(ctx, &up)
		resp.Results = append(resp.Results, result)
		if result.ParseStatus == string(model.ParseSuccess) {
			resp.Success++
		} else {
			resp.Failed++
		}
	}
	return resp, nil
}

// ScanPath 扫描本地路径，将文件上传并入库
func (s *BookFileService) ScanPath(ctx context.Context, path string, userID uint64) (*dto.ScanPathResponse, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %w", err)
	}
	if !info.IsDir() {
		return nil, errors.New("指定的路径不是目录")
	}

	var files []string
	err = filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		ext := strings.ToLower(filepath.Ext(p))
		if ValidateFileType(ext) {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("扫描目录失败: %w", err)
	}

	resp := &dto.ScanPathResponse{Results: make([]dto.ScanResult, 0)}
	for _, fpath := range files {
		result := s.scanLocalFile(ctx, fpath, userID)
		resp.Results = append(resp.Results, result)
		if result.ParseStatus == string(model.ParseSuccess) {
			resp.Imported++
		} else {
			resp.Failed++
		}
	}
	resp.Total = len(files)
	return resp, nil
}

// scanLocalFile 扫描本地文件，上传并入库
func (s *BookFileService) scanLocalFile(ctx context.Context, fpath string, userID uint64) dto.ScanResult {
	originalName := filepath.Base(fpath)
	fileInfo, err := os.Stat(fpath)
	if err != nil {
		return dto.ScanResult{OriginalName: originalName, ParseStatus: string(model.ParseFailed)}
	}
	fileSize := uint64(fileInfo.Size())

	data, err := os.ReadFile(fpath)
	if err != nil {
		return dto.ScanResult{OriginalName: originalName, FileSize: fileSize, ParseStatus: string(model.ParseFailed)}
	}

	md5Hash := fmt.Sprintf("%x", md5.Sum(data))

	// 写入存储
	ext := strings.ToLower(filepath.Ext(originalName))
	relPath := filepath.Join(StorageBaseDir, fmt.Sprintf("%s_%s", md5Hash[:8], originalName))
	absPath := filepath.Join(".", relPath)
	_ = os.MkdirAll(filepath.Dir(absPath), 0755)
	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return dto.ScanResult{OriginalName: originalName, FileSize: fileSize, ParseStatus: string(model.ParseFailed)}
	}

	// 提取元数据
	title, author := extractMetaFromContent(data, originalName)

	// 先应用入库过滤规则，再在过滤后的数据上解析章节
	// 确保字节偏移与存储到文件的内容一致
	filterRules, _ := s.filterRuleRepo.ListByStage(ctx, model.FilterStageInput)
	filter := NewContentFilter(filterRules)
	filteredData := applyFilterToContent(data, filter)

	rules, _ := s.chapterRuleRepo.ListEffective(ctx, 0)
	parser := NewChapterParser(rules)
	parseResult := parser.Parse(filteredData)

	// 匹配或创建 Book
	var bookID uint64
	var chapterCount uint32
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		book, findErr := s.bookRepo.FindByTitleAndAuthor(ctx, title, author)
		if findErr != nil {
			book = &model.Book{
				Title:        title,
				Author:       author,
				SerialStatus: model.SerialOngoing,
				Visibility:   model.VisibilityPublic,
				Status:       model.BookReviewing,
			}
			book.CreateBy = &userID
			book.UpdateBy = &userID
			if err := tx.Create(book).Error; err != nil {
				return err
			}
		}

		// 写内容文件
		contentRelPath := filepath.Join(StorageBaseDir, fmt.Sprintf("content_%d.txt", book.ID))
		contentAbsPath := filepath.Join(".", contentRelPath)
		_ = os.MkdirAll(filepath.Dir(contentAbsPath), 0755)
		if err := os.WriteFile(contentAbsPath, filteredData, 0644); err != nil {
			return err
		}

		// 创建 BookFile
		contentMD5 := fmt.Sprintf("%x", md5.Sum(filteredData))
		charset := DetectCharset(filteredData)
		bf := &model.BookFile{
			BookID:         book.ID,
			OriginalName:   originalName,
			SourceType:     model.FileSourceLocalScan,
			SourceFormat:   &ext,
			ContentPath:    &contentRelPath,
			ContentSize:    uint64(len(filteredData)),
			ContentMD5:     &contentMD5,
			ContentCharset: charset,
			ContentVersion: 1,
			ChapterCount:   uint32(len(parseResult.Chapters)),
			IsPrimary:      true,
			FileStatus:     model.FileSuccess,
		}
		bf.CreateBy = &userID
		bf.UpdateBy = &userID
		if err := tx.Create(bf).Error; err != nil {
			return err
		}

		// 创建章节索引
		chapters := make([]model.BookChapter, len(parseResult.Chapters))
		for i, seg := range parseResult.Chapters {
			chapters[i] = model.BookChapter{
				BookID:     book.ID,
				FileID:     bf.ID,
				ChapterNo:  uint32(i + 1),
				Title:      seg.Title,
				ByteOffset: seg.ByteOffset,
				ByteLength: seg.ByteLength,
				WordCount:  seg.WordCount,
				Status:     model.ChapterPublished,
			}
		}
		if err := tx.Create(&chapters).Error; err != nil {
			return err
		}

		// 更新书籍统计
		book.PrimaryFileID = &bf.ID
		book.TotalChapters = uint32(len(parseResult.Chapters))
		book.TotalWords = sumWords(parseResult.Chapters)
		book.UpdateBy = &userID
		if err := tx.Save(book).Error; err != nil {
			return err
		}

		bookID = book.ID
		chapterCount = uint32(len(parseResult.Chapters))
		return nil
	})

	parseStatus := string(model.ParseSuccess)
	var parseMsg *string
	if err != nil {
		parseStatus = string(model.ParseFailed)
		msg := err.Error()
		parseMsg = &msg
	}

	return dto.ScanResult{
		OriginalName: originalName,
		FileSize:     fileSize,
		ParseStatus:  parseStatus,
		ParseMessage: parseMsg,
		BookID:       &bookID,
		ChapterCount: &chapterCount,
	}
}

// ScanByID 扫描单个上传任务
func (s *BookFileService) ScanByID(ctx context.Context, uploadID uint64) (*dto.ScanResult, error) {
	up, err := s.uploadRepo.GetByID(ctx, uploadID)
	if err != nil {
		return nil, ErrUploadNotFound
	}
	result := s.scanSingle(ctx, up)
	return &result, nil
}

// scanSingle 解析单个文件并入库
func (s *BookFileService) scanSingle(ctx context.Context, up *model.BookUpload) dto.ScanResult {
	result := dto.ScanResult{
		UploadID:     up.ID,
		OriginalName: up.OriginalName,
		FileSize:     up.FileSize,
		ParseStatus:  string(model.ParseProcessing),
	}

	// 标记处理中
	up.ParseStatus = model.ParseProcessing
	_ = s.uploadRepo.Update(ctx, up)

	// 读取文件
	data, err := os.ReadFile(up.FilePath)
	if err != nil {
		failMsg := fmt.Sprintf("读取文件失败: %v", err)
		up.ParseStatus = model.ParseFailed
		up.ParseMessage = &failMsg
		_ = s.uploadRepo.Update(ctx, up)
		result.ParseStatus = string(model.ParseFailed)
		result.ParseMessage = &failMsg
		return result
	}

	// 获取章节识别规则
	rules, _ := s.chapterRuleRepo.ListEffective(ctx, 0) // bookID=0 取全局规则

	// 先应用入库过滤规则，再在过滤后的数据上解析章节
	// 确保字节偏移与存储到文件的内容一致
	filterRules, _ := s.filterRuleRepo.ListByStage(ctx, model.FilterStageInput)
	filter := NewContentFilter(filterRules)
	filteredData := applyFilterToContent(data, filter)

	parser := NewChapterParser(rules)
	parseResult := parser.Parse(filteredData)

	// 事务入库
	var bookID uint64
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建 Book 记录
		title, author := extractMetaFromContent(data, up.OriginalName)
		book := &model.Book{
			Title:         title,
			Author:        author,
			SerialStatus:  model.SerialOngoing,
			Visibility:    model.VisibilityPublic,
			TotalChapters: uint32(len(parseResult.Chapters)),
			TotalWords:    sumWords(parseResult.Chapters),
			Status:        model.BookReviewing,
		}
		if err := tx.Create(book).Error; err != nil {
			return err
		}
		bookID = book.ID

		// 写入解析后的内容文件
		contentRelPath := filepath.Join(StorageBaseDir, fmt.Sprintf("content_%d.txt", bookID))
		contentAbsPath := filepath.Join(".", contentRelPath)
		if err := os.MkdirAll(filepath.Dir(contentAbsPath), 0755); err != nil {
			return fmt.Errorf("创建内容目录失败: %w", err)
		}
		if err := os.WriteFile(contentAbsPath, filteredData, 0644); err != nil {
			return fmt.Errorf("写入内容文件失败: %w", err)
		}

		// 创建 BookFile 记录
		contentMD5 := fmt.Sprintf("%x", md5.Sum(filteredData))
		charset := DetectCharset(filteredData)
		bf := &model.BookFile{
			BookID:         bookID,
			OriginalName:   up.OriginalName,
			SourceType:     model.FileSourceUserUpload,
			SourceFormat:   up.SourceFormat,
			SourceFileURL:  &up.FilePath,
			ContentPath:    &contentRelPath,
			ContentSize:    uint64(len(filteredData)),
			ContentMD5:     &contentMD5,
			ContentCharset: charset,
			ContentVersion: 1,
			ChapterCount:   uint32(len(parseResult.Chapters)),
			IsPrimary:      true,
			FileStatus:     model.FileSuccess,
		}
		if err := tx.Create(bf).Error; err != nil {
			return err
		}

		// 创建章节索引
		chapters := make([]model.BookChapter, len(parseResult.Chapters))
		for i, seg := range parseResult.Chapters {
			chapters[i] = model.BookChapter{
				BookID:     bookID,
				FileID:     bf.ID,
				ChapterNo:  uint32(i + 1),
				Title:      seg.Title,
				ByteOffset: seg.ByteOffset,
				ByteLength: seg.ByteLength,
				WordCount:  seg.WordCount,
				Status:     model.ChapterPublished,
			}
		}
		if err := tx.Create(&chapters).Error; err != nil {
			return err
		}

		// 更新上传记录
		chapterCount := uint32(len(parseResult.Chapters))
		up.BookID = &bookID
		up.ParseStatus = model.ParseSuccess
		up.ChapterCount = &chapterCount
		if err := tx.Save(up).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		failMsg := fmt.Sprintf("入库失败: %v", err)
		up.ParseStatus = model.ParseFailed
		up.ParseMessage = &failMsg
		_ = s.uploadRepo.Update(ctx, up)
		result.ParseStatus = string(model.ParseFailed)
		result.ParseMessage = &failMsg
		return result
	}

	result.ParseStatus = string(model.ParseSuccess)
	result.BookID = &bookID
	cc := uint32(len(parseResult.Chapters))
	result.ChapterCount = &cc
	return result
}

// ==================== 章节内容读取 ====================

// ReParseChapters 重新识别章节
// 读取已入库的内容文件，重新解析章节索引并替换
func (s *BookFileService) ReParseChapters(ctx context.Context, req *dto.ReParseRequest) (*dto.ReParseResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, req.BookID)
	if err != nil {
		return nil, fmt.Errorf("书籍不存在: %w", err)
	}
	if book.PrimaryFileID == nil {
		return nil, errors.New("该书没有主文件，无法重新识别章节")
	}
	file, err := s.bookFileRepo.GetByID(ctx, *book.PrimaryFileID)
	if err != nil {
		return nil, ErrBookFileNotFound
	}
	if file.ContentPath == nil {
		return nil, errors.New("文件路径为空")
	}
	data, err := os.ReadFile(*file.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("读取内容文件失败: %w", err)
	}

	// 获取章节识别规则
	rules, _ := s.chapterRuleRepo.ListEffective(ctx, 0)

	// 在已过滤的内容上重新解析
	parser := NewChapterParser(rules)
	parseResult := parser.Parse(data)

	oldCount := book.TotalChapters
	newCount := uint32(len(parseResult.Chapters))

	// 事务替换章节索引
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 硬删除旧章节（唯一索引 uk_book_file_chapter 不含 deleted_at）
		//    必须用 Unscoped() 彻底删除，否则软删除的记录会阻塞新插入
		if err := tx.Unscoped().Where("book_id = ?", book.ID).Delete(&model.BookChapter{}).Error; err != nil {
			return fmt.Errorf("删除旧章节失败: %w", err)
		}
		// 2. 创建新章节
		chapters := make([]model.BookChapter, newCount)
		for i, seg := range parseResult.Chapters {
			chapters[i] = model.BookChapter{
				BookID:     book.ID,
				FileID:     *book.PrimaryFileID,
				ChapterNo:  uint32(i + 1),
				Title:      seg.Title,
				ByteOffset: seg.ByteOffset,
				ByteLength: seg.ByteLength,
				WordCount:  seg.WordCount,
				Status:     model.ChapterPublished,
			}
		}
		if err := tx.Create(&chapters).Error; err != nil {
			return fmt.Errorf("创建新章节失败: %w", err)
		}
		// 3. 更新书籍统计
		totalWords := sumWords(parseResult.Chapters)
		if err := tx.Model(&model.Book{}).Where("id = ?", book.ID).Updates(map[string]interface{}{
			"total_chapters": newCount,
			"total_words":    totalWords,
		}).Error; err != nil {
			return fmt.Errorf("更新书籍统计失败: %w", err)
		}
		// 4. 更新文件记录的章节数
		if err := tx.Model(&model.BookFile{}).Where("id = ?", file.ID).Update("chapter_count", newCount).Error; err != nil {
			return fmt.Errorf("更新文件章节数失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.ReParseResponse{
		BookID:     book.ID,
		BookTitle:  book.Title,
		OldCount:   oldCount,
		NewCount:   newCount,
		TotalWords: sumWords(parseResult.Chapters),
	}, nil
}

// GetChapterContent 读取指定章节的文本内容
func (s *BookFileService) GetChapterContent(ctx context.Context, bookID uint64, chapterNo uint32) (*dto.ChapterContentResponse, error) {
	chapter, err := s.chapterRepo.GetByBookAndNo(ctx, bookID, chapterNo)
	if err != nil {
		return nil, ErrChapterNotFound
	}
	file, err := s.bookFileRepo.GetByID(ctx, chapter.FileID)
	if err != nil {
		return nil, ErrBookFileNotFound
	}
	if file.ContentPath == nil {
		return nil, errors.New("文件路径为空")
	}
	data, err := os.ReadFile(*file.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("读取章节文件失败: %w", err)
	}
	content := string(data[chapter.ByteOffset : chapter.ByteOffset+uint64(chapter.ByteLength)])

	// 应用出库过滤规则
	filterRules, _ := s.filterRuleRepo.ListByStage(ctx, model.FilterStageOutput)
	if len(filterRules) > 0 {
		filter := NewContentFilter(filterRules)
		fr := filter.Filter(content)
		if fr.Blocked {
			return nil, errors.New(fr.MatchDesc)
		}
		content = fr.Content
	}

	return &dto.ChapterContentResponse{
		BookChapter: *chapter,
		Content:     content,
	}, nil
}

// ==================== 章节识别规则 CRUD ====================

func (s *BookFileService) CreateChapterRule(ctx context.Context, req *dto.ChapterRuleRequest, userID uint64) (*model.BookChapterRule, error) {
	m := &model.BookChapterRule{
		RuleName:      req.RuleName,
		ScopeType:     model.RuleScopeType(req.ScopeType),
		BookID:        req.BookID,
		Pattern:       req.Pattern,
		TitleGroup:    req.TitleGroup,
		MinChapterLen: req.MinChapterLen,
		MaxChapterLen: req.MaxChapterLen,
		Priority:      req.Priority,
		Description:   req.Description,
		Status:        model.EnableStatus(req.Status),
	}
	if m.Status == "" {
		m.Status = model.StatusEnabled
	}
	if m.MinChapterLen == 0 {
		m.MinChapterLen = DefaultChapterLen
	}
	if m.MaxChapterLen == 0 {
		m.MaxChapterLen = MaxChapterLen
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.chapterRuleRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookFileService) UpdateChapterRule(ctx context.Context, id uint64, req *dto.ChapterRuleRequest, userID uint64) (*model.BookChapterRule, error) {
	m, err := s.chapterRuleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrRuleNotFound
	}
	m.RuleName = req.RuleName
	m.ScopeType = model.RuleScopeType(req.ScopeType)
	m.BookID = req.BookID
	m.Pattern = req.Pattern
	m.TitleGroup = req.TitleGroup
	m.MinChapterLen = req.MinChapterLen
	m.MaxChapterLen = req.MaxChapterLen
	m.Priority = req.Priority
	m.Description = req.Description
	if req.Status != "" {
		m.Status = model.EnableStatus(req.Status)
	}
	m.UpdateBy = &userID
	if err := s.chapterRuleRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookFileService) DeleteChapterRule(ctx context.Context, id uint64) error {
	if _, err := s.chapterRuleRepo.GetByID(ctx, id); err != nil {
		return ErrRuleNotFound
	}
	return s.chapterRuleRepo.Delete(ctx, id)
}

func (s *BookFileService) GetChapterRuleByID(ctx context.Context, id uint64) (*model.BookChapterRule, error) {
	return s.chapterRuleRepo.GetByID(ctx, id)
}

func (s *BookFileService) PageChapterRule(ctx context.Context, req *dto.ChapterRuleSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.chapterRuleRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	records := make([]dto.ChapterRuleResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.ChapterRuleResponse{BookChapterRule: r}
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

// ==================== 内容净化规则 CRUD ====================

func (s *BookFileService) CreateFilterRule(ctx context.Context, req *dto.FilterRuleRequest, userID uint64) (*model.BookContentFilterRule, error) {
	m := &model.BookContentFilterRule{
		RuleName:    req.RuleName,
		MatchType:   model.FilterMatchType(req.MatchType),
		Pattern:     req.Pattern,
		Action:      model.FilterAction(req.Action),
		Replacement: req.Replacement,
		ApplyStage:  model.FilterApplyStage(req.ApplyStage),
		Category:    req.Category,
		Severity:    model.FilterSeverity(req.Severity),
		Description: req.Description,
		Status:      model.EnableStatus(req.Status),
	}
	if m.Replacement == "" {
		m.Replacement = "***"
	}
	if m.Severity == "" {
		m.Severity = model.FilterLow
	}
	if m.Status == "" {
		m.Status = model.StatusEnabled
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.filterRuleRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookFileService) UpdateFilterRule(ctx context.Context, id uint64, req *dto.FilterRuleRequest, userID uint64) (*model.BookContentFilterRule, error) {
	m, err := s.filterRuleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrFilterRuleNotFound
	}
	m.RuleName = req.RuleName
	m.MatchType = model.FilterMatchType(req.MatchType)
	m.Pattern = req.Pattern
	m.Action = model.FilterAction(req.Action)
	m.Replacement = req.Replacement
	m.ApplyStage = model.FilterApplyStage(req.ApplyStage)
	m.Category = req.Category
	m.Severity = model.FilterSeverity(req.Severity)
	m.Description = req.Description
	if req.Status != "" {
		m.Status = model.EnableStatus(req.Status)
	}
	m.UpdateBy = &userID
	if err := s.filterRuleRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookFileService) DeleteFilterRule(ctx context.Context, id uint64) error {
	if _, err := s.filterRuleRepo.GetByID(ctx, id); err != nil {
		return ErrFilterRuleNotFound
	}
	return s.filterRuleRepo.Delete(ctx, id)
}

func (s *BookFileService) GetFilterRuleByID(ctx context.Context, id uint64) (*model.BookContentFilterRule, error) {
	return s.filterRuleRepo.GetByID(ctx, id)
}

func (s *BookFileService) PageFilterRule(ctx context.Context, req *dto.FilterRuleSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.filterRuleRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	records := make([]dto.FilterRuleResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.FilterRuleResponse{BookContentFilterRule: r}
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

// ==================== 上传/文件查询 ====================

func (s *BookFileService) PageUpload(ctx context.Context, req *dto.UploadSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.uploadRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	records := make([]dto.UploadResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.UploadResponse{BookUpload: r}
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

func (s *BookFileService) PageFile(ctx context.Context, req *dto.FileSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.bookFileRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	records := make([]dto.FileResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.FileResponse{BookFile: r}
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

func (s *BookFileService) PageChapter(ctx context.Context, req *dto.ChapterSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.chapterRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	records := make([]dto.ChapterResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.ChapterResponse{BookChapter: r}
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

// ListChapter 不分页章节列表
func (s *BookFileService) ListChapter(ctx context.Context, req *dto.ChapterListRequest) ([]dto.ChapterResponse, error) {
	rows, err := s.chapterRepo.ListByBookID(ctx, req.BookID)
	if err != nil {
		return nil, err
	}
	records := make([]dto.ChapterResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.ChapterResponse{BookChapter: r}
	}
	return records, nil
}

// ==================== 辅助函数 ====================

// tryDecodeToUTF8 尝试将 GBK 编码数据转为 UTF-8，检测失败则返回原数据
func tryDecodeToUTF8(data []byte) []byte {
	// 如果已经是有效 UTF-8，直接返回
	if utf8.Valid(data) {
		return data
	}
	decoded, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return data
	}
	return decoded
}

func sumWords(chapters []ChapterSegment) uint32 {
	var total uint32
	for _, c := range chapters {
		total += c.WordCount
	}
	return total
}

// applyFilterToContent 对整本书的字节数据应用过滤规则
func applyFilterToContent(data []byte, filter *ContentFilter) []byte {
	lines := strings.Split(string(data), "\n")
	var result []string
	for _, line := range lines {
		fr := filter.Filter(line)
		if fr.Blocked {
			result = append(result, "")
			continue
		}
		result = append(result, fr.Content)
	}
	return []byte(strings.Join(result, "\n"))
}
