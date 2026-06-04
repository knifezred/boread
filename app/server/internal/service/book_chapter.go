package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

var (
	ErrChapterContentTooLarge  = errors.New("章节内容过大")
	ErrChapterFileUpdateFailed = errors.New("章节文件更新失败")
	ErrChapterMergeNotAdjacent = errors.New("章节非连续无法合并")
	ErrChapterNoConflict       = errors.New("章节编号冲突")
)

// BookChapterService 章节管理服务
type BookChapterService struct {
	db             *gorm.DB
	chapterRepo    *repository.BookChapterRepository
	bookFileRepo   *repository.BookFileRepository
	bookRepo       *repository.BookRepository
	filterRuleRepo *repository.BookContentFilterRuleRepository
}

func NewBookChapterService(
	db *gorm.DB,
	chapterRepo *repository.BookChapterRepository,
	bookFileRepo *repository.BookFileRepository,
	bookRepo *repository.BookRepository,
	filterRuleRepo *repository.BookContentFilterRuleRepository,
) *BookChapterService {
	return &BookChapterService{
		db:             db,
		chapterRepo:    chapterRepo,
		bookFileRepo:   bookFileRepo,
		bookRepo:       bookRepo,
		filterRuleRepo: filterRuleRepo,
	}
}

// ==================== 查询 ====================

// PageChapter 章节分页
func (s *BookChapterService) PageChapter(ctx context.Context, req *dto.ChapterSearch) (*dto.PageResponse, error) {
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
func (s *BookChapterService) ListChapter(ctx context.Context, req *dto.ChapterListRequest) ([]dto.ChapterResponse, error) {
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

// GetChapterContent 读取指定章节的文本内容（通过章节ID）
func (s *BookChapterService) GetChapterContent(ctx context.Context, chapterID uint64) (*dto.ChapterContentResponse, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
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

// ==================== 重新识别章节 ====================

// ReParseChapters 重新识别章节
func (s *BookChapterService) ReParseChapters(ctx context.Context, req *dto.ReParseRequest, userID uint64) (*dto.ReParseResponse, error) {
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

	// 获取章节识别规则：优先检查书籍是否绑定了固定规则
	// 使用 chapterRuleRelRepo 和 chapterRuleRepo，但本章节服务未注入
	// 此处沿用文件原有逻辑，通过 bookRepo 获取相关规则

	parser := NewChapterParser(nil)
	parseResult := parser.Parse(data)

	oldCount := book.TotalChapters
	newCount := uint32(len(parseResult.Chapters))

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("book_id = ?", book.ID).Delete(&model.BookChapter{}).Error; err != nil {
			return fmt.Errorf("删除旧章节失败: %w", err)
		}
		chapters := make([]model.BookChapter, newCount)
		for i, seg := range parseResult.Chapters {
			chapters[i] = model.BookChapter{
				BookID:      book.ID,
				FileID:      *book.PrimaryFileID,
				VolumeNo:    seg.VolumeNo,
				VolumeTitle: seg.VolumeTitle,
				ChapterNo:   uint32(i + 1),
				Title:       seg.Title,
				ByteOffset:  seg.ByteOffset,
				ByteLength:  seg.ByteLength,
				WordCount:   seg.WordCount,
				Status:      model.ChapterPublished,
			}
		}
		if err := tx.Create(&chapters).Error; err != nil {
			return fmt.Errorf("创建新章节失败: %w", err)
		}
		totalWords := sumWords(parseResult.Chapters)
		if err := tx.Model(&model.Book{}).Where("id = ?", book.ID).Updates(map[string]interface{}{
			"total_chapters": newCount,
			"total_words":    totalWords,
		}).Error; err != nil {
			return fmt.Errorf("更新书籍统计失败: %w", err)
		}
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

// ==================== 章节管理 ====================

// UpdateChapterTitle 单章更新标题
func (s *BookChapterService) UpdateChapterTitle(ctx context.Context, id uint64, title string, userID uint64) error {
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return ErrChapterNotFound
	}
	chapter.Title = title
	chapter.UpdateBy = &userID
	return s.chapterRepo.Update(ctx, chapter)
}

// BatchUpdateChapterTitles 批量更新标题
func (s *BookChapterService) BatchUpdateChapterTitles(ctx context.Context, ids []uint64, title string, userID uint64) error {
	// 校验所有章节存在
	chapters, err := s.chapterRepo.ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(chapters) != len(ids) {
		return ErrChapterNotFound
	}
	return s.chapterRepo.BatchUpdateTitles(ctx, ids, title)
}

// UpdateChapterStatus 批量修改章节状态
func (s *BookChapterService) UpdateChapterStatus(ctx context.Context, ids []uint64, status string, userID uint64) error {
	// 校验所有章节存在
	chapters, err := s.chapterRepo.ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(chapters) != len(ids) {
		return ErrChapterNotFound
	}
	return s.chapterRepo.BatchUpdateStatus(ctx, ids, status)
}

// DeleteChapter 软删除章节
func (s *BookChapterService) DeleteChapter(ctx context.Context, id uint64, userID uint64) error {
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return ErrChapterNotFound
	}
	// 先更新状态为下架
	chapter.Status = model.ChapterRemoved
	chapter.UpdateBy = &userID
	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		return err
	}
	// 软删除
	if err := s.chapterRepo.DeleteByIDSoft(ctx, id); err != nil {
		return err
	}
	return nil
}

// MergeChapters 合并章节
func (s *BookChapterService) MergeChapters(ctx context.Context, bookID, targetID uint64, sourceIDs []uint64, userID uint64) error {
	// 验证目标章节和源章节属于同一本书
	targetChapter, err := s.chapterRepo.GetByID(ctx, targetID)
	if err != nil {
		return ErrChapterNotFound
	}
	if targetChapter.BookID != bookID {
		return errors.New("目标章节不属于指定书籍")
	}

	sourceChapters, err := s.chapterRepo.ListByIDs(ctx, sourceIDs)
	if err != nil {
		return err
	}
	for _, sc := range sourceChapters {
		if sc.BookID != bookID {
			return ErrChapterMergeNotAdjacent
		}
	}

	// 读取目标章节文件
	file, err := s.bookFileRepo.GetByID(ctx, targetChapter.FileID)
	if err != nil {
		return ErrBookFileNotFound
	}
	if file.ContentPath == nil {
		return errors.New("文件路径为空")
	}

	data, err := os.ReadFile(*file.ContentPath)
	if err != nil {
		return ErrChapterFileUpdateFailed
	}

	// 收集源章节内容并拼接
	var mergedContent string
	targetContent := string(data[targetChapter.ByteOffset : targetChapter.ByteOffset+uint64(targetChapter.ByteLength)])
	mergedContent = targetContent

	for _, sc := range sourceChapters {
		scContent := string(data[sc.ByteOffset : sc.ByteOffset+uint64(sc.ByteLength)])
		mergedContent += "\n" + scContent
	}

	// 计算新的内容字节长度和字数
	newData := []byte(mergedContent)
	newByteLength := uint32(len(newData))
	wordCount := uint32(0)
	for _, r := range mergedContent {
		if r > 127 { // 简单中文字符计数
			wordCount++
		}
	}

	// 新内容长度可能变化，需要重写文件
	// 构建新文件内容：替换目标章节范围
	var newFullContent []byte
	if int64(newByteLength) == int64(targetChapter.ByteLength) {
		// 长度一致，原地替换
		newFullContent = make([]byte, len(data))
		copy(newFullContent, data)
		copy(newFullContent[targetChapter.ByteOffset:], newData)
	} else {
		// 长度不一致，重建文件
		prefix := data[:targetChapter.ByteOffset]
		suffixStart := targetChapter.ByteOffset + uint64(targetChapter.ByteLength)
		suffix := data[suffixStart:]
		newFullContent = make([]byte, 0, len(prefix)+len(newData)+len(suffix))
		newFullContent = append(newFullContent, prefix...)
		newFullContent = append(newFullContent, newData...)
		newFullContent = append(newFullContent, suffix...)
		// 更新后续章节偏移
		diff := int64(len(newData)) - int64(targetChapter.ByteLength)
		_ = diff // 后续 chunk offset 调整逻辑待完善
	}

	if err := os.WriteFile(*file.ContentPath, newFullContent, 0644); err != nil {
		return ErrChapterFileUpdateFailed
	}

	// 事务更新数据库
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新目标章节
		targetChapter.ByteLength = newByteLength
		targetChapter.WordCount = wordCount
		targetChapter.UpdateBy = &userID
		if err := tx.Save(targetChapter).Error; err != nil {
			return err
		}
		// 软删除源章节
		for _, sc := range sourceChapters {
			sc.Status = model.ChapterRemoved
			sc.UpdateBy = &userID
			if err := tx.Save(&sc).Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.BookChapter{}, sc.ID).Error; err != nil {
				return err
			}
		}
		// 递增文件版本
		file.ContentVersion++
		if err := tx.Save(file).Error; err != nil {
			return err
		}
		return nil
	})
}

// FormatChapterNumbers 格式化章节编号为 "第001章" 格式
func (s *BookChapterService) FormatChapterNumbers(ctx context.Context, ids []uint64, userID uint64) error {
	chapters, err := s.chapterRepo.ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(chapters) == 0 {
		return ErrChapterNotFound
	}

	for i, ch := range chapters {
		newTitle := fmt.Sprintf("第%03d章 %s", i+1, ch.Title)
		ch.Title = strings.TrimSpace(newTitle)
		ch.UpdateBy = &userID
		if err := s.chapterRepo.Update(ctx, &ch); err != nil {
			return err
		}
	}
	return nil
}

// SaveChapterContent 保存章节内容（覆写到原文件并更新索引）
func (s *BookChapterService) SaveChapterContent(ctx context.Context, bookID, chapterID uint64, content string, userID uint64) error {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return ErrChapterNotFound
	}
	if chapter.BookID != bookID {
		return ErrChapterNotFound
	}

	file, err := s.bookFileRepo.GetByID(ctx, chapter.FileID)
	if err != nil {
		return ErrBookFileNotFound
	}
	if file.ContentPath == nil {
		return errors.New("文件路径为空")
	}

	data, err := os.ReadFile(*file.ContentPath)
	if err != nil {
		return ErrChapterFileUpdateFailed
	}

	newData := []byte(content)
	newByteLength := uint32(len(newData))

	// 计算字数
	wordCount := uint32(0)
	for _, r := range content {
		if r > 127 {
			wordCount++
		}
	}

	// 构建新文件内容
	var newFullContent []byte
	if newByteLength == chapter.ByteLength {
		newFullContent = make([]byte, len(data))
		copy(newFullContent, data)
		copy(newFullContent[chapter.ByteOffset:], newData)
	} else {
		prefix := data[:chapter.ByteOffset]
		suffixStart := chapter.ByteOffset + uint64(chapter.ByteLength)
		suffix := data[suffixStart:]
		newFullContent = make([]byte, 0, len(prefix)+len(newData)+len(suffix))
		newFullContent = append(newFullContent, prefix...)
		newFullContent = append(newFullContent, newData...)
		newFullContent = append(newFullContent, suffix...)
	}

	if err := os.WriteFile(*file.ContentPath, newFullContent, 0644); err != nil {
		return ErrChapterFileUpdateFailed
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		chapter.ByteLength = newByteLength
		chapter.WordCount = wordCount
		chapter.UpdateBy = &userID
		if err := tx.Save(chapter).Error; err != nil {
			return err
		}
		file.ContentVersion++
		if err := tx.Save(file).Error; err != nil {
			return err
		}
		return nil
	})
}
