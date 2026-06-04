package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

// ==================== BookUploadRepository ====================

type BookUploadRepository struct {
	db *gorm.DB
}

func NewBookUploadRepository(db *gorm.DB) *BookUploadRepository {
	return &BookUploadRepository{db: db}
}

func (r *BookUploadRepository) Create(ctx context.Context, m *model.BookUpload) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookUploadRepository) Update(ctx context.Context, m *model.BookUpload) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookUploadRepository) GetByID(ctx context.Context, id uint64) (*model.BookUpload, error) {
	var m model.BookUpload
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookUploadRepository) Page(ctx context.Context, req *dto.UploadSearch) ([]model.BookUpload, int64, error) {
	var rows []model.BookUpload
	tx := r.db.WithContext(ctx).Model(&model.BookUpload{})
	if req.OriginalName != "" {
		tx = tx.Where("original_name LIKE ?", "%"+req.OriginalName+"%")
	}
	if req.ParseStatus != "" {
		tx = tx.Where("parse_status = ?", req.ParseStatus)
	}
	if req.BookID != nil {
		tx = tx.Where("book_id = ?", *req.BookID)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("id DESC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *BookUploadRepository) ListByParseStatus(ctx context.Context, status model.ParseStatus, limit int) ([]model.BookUpload, error) {
	var rows []model.BookUpload
	err := r.db.WithContext(ctx).Model(&model.BookUpload{}).
		Where("parse_status = ?", status).
		Order("id ASC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

// ==================== BookFileRepository ====================

type BookFileRepository struct {
	db *gorm.DB
}

func NewBookFileRepository(db *gorm.DB) *BookFileRepository {
	return &BookFileRepository{db: db}
}

func (r *BookFileRepository) Create(ctx context.Context, m *model.BookFile) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookFileRepository) Update(ctx context.Context, m *model.BookFile) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookFileRepository) GetByID(ctx context.Context, id uint64) (*model.BookFile, error) {
	var m model.BookFile
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookFileRepository) GetByBookID(ctx context.Context, bookID uint64) ([]model.BookFile, error) {
	var rows []model.BookFile
	err := r.db.WithContext(ctx).Model(&model.BookFile{}).
		Where("book_id = ?", bookID).
		Order("id ASC").
		Find(&rows).Error
	return rows, err
}

func (r *BookFileRepository) GetPrimaryByBookID(ctx context.Context, bookID uint64) (*model.BookFile, error) {
	var m model.BookFile
	err := r.db.WithContext(ctx).Model(&model.BookFile{}).
		Where("book_id = ? AND is_primary = 1", bookID).
		First(&m).Error
	return &m, err
}

func (r *BookFileRepository) Page(ctx context.Context, req *dto.FileSearch) ([]model.BookFile, int64, error) {
	var rows []model.BookFile
	tx := r.db.WithContext(ctx).Model(&model.BookFile{})
	if req.BookID != nil {
		tx = tx.Where("book_id = ?", *req.BookID)
	}
	if req.FileStatus != "" {
		tx = tx.Where("file_status = ?", req.FileStatus)
	}
	if req.SourceType != "" {
		tx = tx.Where("source_type = ?", req.SourceType)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("id DESC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ==================== BookChapterRepository ====================

type BookChapterRepository struct {
	db *gorm.DB
}

func NewBookChapterRepository(db *gorm.DB) *BookChapterRepository {
	return &BookChapterRepository{db: db}
}

func (r *BookChapterRepository) BatchCreate(ctx context.Context, chapters []model.BookChapter) error {
	if len(chapters) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&chapters).Error
}

func (r *BookChapterRepository) DeleteByFileID(ctx context.Context, fileID uint64) error {
	return r.db.WithContext(ctx).Where("file_id = ?", fileID).Delete(&model.BookChapter{}).Error
}

func (r *BookChapterRepository) GetLatestByBookID(ctx context.Context, bookID uint64) (*model.BookChapter, error) {
	var m model.BookChapter
	err := r.db.WithContext(ctx).Model(&model.BookChapter{}).
		Where("book_id = ?", bookID).
		Order("chapter_no DESC").
		Limit(1).
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	if m.ID == 0 {
		return nil, nil
	}
	return &m, nil
}

func (r *BookChapterRepository) GetByID(ctx context.Context, id uint64) (*model.BookChapter, error) {
	var m model.BookChapter
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookChapterRepository) DeleteByBookID(ctx context.Context, bookID uint64) error {
	return r.db.WithContext(ctx).Where("book_id = ?", bookID).Delete(&model.BookChapter{}).Error
}

func (r *BookChapterRepository) GetByBookAndNo(ctx context.Context, bookID uint64, chapterNo uint32) (*model.BookChapter, error) {
	var m model.BookChapter
	err := r.db.WithContext(ctx).Where("book_id = ? AND chapter_no = ?", bookID, chapterNo).First(&m).Error
	return &m, err
}

func (r *BookChapterRepository) Page(ctx context.Context, req *dto.ChapterSearch) ([]model.BookChapter, int64, error) {
	var rows []model.BookChapter
	tx := r.db.WithContext(ctx).Model(&model.BookChapter{})
	if req.BookID != nil {
		tx = tx.Where("book_id = ?", *req.BookID)
	}
	if req.FileID != nil {
		tx = tx.Where("file_id = ?", *req.FileID)
	}
	if req.ChapterNo != nil {
		tx = tx.Where("chapter_no = ?", *req.ChapterNo)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("chapter_no ASC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *BookChapterRepository) ListByBookID(ctx context.Context, bookID uint64) ([]model.BookChapter, error) {
	var rows []model.BookChapter
	err := r.db.WithContext(ctx).Model(&model.BookChapter{}).
		Where("book_id = ?", bookID).
		Order("chapter_no ASC").
		Find(&rows).Error
	return rows, err
}

func (r *BookChapterRepository) GetMaxChapterNo(ctx context.Context, bookID, fileID uint64) (uint32, error) {
	var maxNo uint32
	err := r.db.WithContext(ctx).Model(&model.BookChapter{}).
		Where("book_id = ? AND file_id = ?", bookID, fileID).
		Select("COALESCE(MAX(chapter_no), 0)").
		Scan(&maxNo).Error
	return maxNo, err
}

// ==================== BookChapterRuleRepository ====================

type BookChapterRuleRepository struct {
	db *gorm.DB
}

func NewBookChapterRuleRepository(db *gorm.DB) *BookChapterRuleRepository {
	return &BookChapterRuleRepository{db: db}
}

func (r *BookChapterRuleRepository) Create(ctx context.Context, m *model.BookChapterRule) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookChapterRuleRepository) Update(ctx context.Context, m *model.BookChapterRule) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookChapterRuleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookChapterRule{}, id).Error
}

func (r *BookChapterRuleRepository) GetByID(ctx context.Context, id uint64) (*model.BookChapterRule, error) {
	var m model.BookChapterRule
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookChapterRuleRepository) Page(ctx context.Context, req *dto.ChapterRuleSearch) ([]model.BookChapterRule, int64, error) {
	var rows []model.BookChapterRule
	tx := r.db.WithContext(ctx).Model(&model.BookChapterRule{})
	if req.RuleName != "" {
		tx = tx.Where("rule_name LIKE ?", "%"+req.RuleName+"%")
	}
	if req.RuleType != "" {
		tx = tx.Where("rule_type = ?", req.RuleType)
	}
	if req.UserID != nil {
		tx = tx.Where("user_id = ?", *req.UserID)
	}
	if req.Status != "" {
		tx = tx.Where("status = ?", req.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("sort_order ASC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *BookChapterRuleRepository) ListEffective(ctx context.Context, userID uint64) ([]model.BookChapterRule, error) {
	var rows []model.BookChapterRule
	// 查询启用的规则：用户自定义规则优先，不足时补充系统默认规则
	// 规则匹配: 用户自定义(user_id=userID) + 系统默认(user_id IS NULL)
	err := r.db.WithContext(ctx).Model(&model.BookChapterRule{}).
		Where("status = ?", model.StatusEnabled).
		Where("(user_id = ?) OR (rule_type = ? AND user_id IS NULL)", userID, model.RuleTypeSystem).
		Order("sort_order ASC").
		Find(&rows).Error
	return rows, err
}

func (r *BookChapterRuleRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.BookChapterRule, error) {
	var rows []model.BookChapterRule
	err := r.db.WithContext(ctx).Model(&model.BookChapterRule{}).
		Where("user_id = ?", userID).
		Order("sort_order ASC").
		Find(&rows).Error
	return rows, err
}

func (r *BookChapterRuleRepository) ListSystemDefaults(ctx context.Context) ([]model.BookChapterRule, error) {
	var rows []model.BookChapterRule
	err := r.db.WithContext(ctx).Model(&model.BookChapterRule{}).
		Where("rule_type = ? AND user_id IS NULL", model.RuleTypeSystem).
		Order("sort_order ASC").
		Find(&rows).Error
	return rows, err
}

// ==================== BookChapterRuleRelRepository ====================

type BookChapterRuleRelRepository struct {
	db *gorm.DB
}

func NewBookChapterRuleRelRepository(db *gorm.DB) *BookChapterRuleRelRepository {
	return &BookChapterRuleRelRepository{db: db}
}

func (r *BookChapterRuleRelRepository) Create(ctx context.Context, m *model.BookChapterRuleRel) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookChapterRuleRelRepository) GetByBookAndReader(ctx context.Context, bookID, readerID uint64) (*model.BookChapterRuleRel, error) {
	var m model.BookChapterRuleRel
	err := r.db.WithContext(ctx).
		Where("book_id = ? AND reader_id = ?", bookID, readerID).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookChapterRuleRelRepository) DeleteByBookAndReader(ctx context.Context, bookID, readerID uint64) error {
	return r.db.WithContext(ctx).
		Where("book_id = ? AND reader_id = ?", bookID, readerID).
		Delete(&model.BookChapterRuleRel{}).Error
}

func (r *BookChapterRuleRelRepository) ListByReader(ctx context.Context, readerID uint64) ([]model.BookChapterRuleRel, error) {
	var rows []model.BookChapterRuleRel
	err := r.db.WithContext(ctx).
		Where("reader_id = ?", readerID).
		Find(&rows).Error
	return rows, err
}

// ==================== BookContentFilterRuleRepository ====================

type BookContentFilterRuleRepository struct {
	db *gorm.DB
}

func NewBookContentFilterRuleRepository(db *gorm.DB) *BookContentFilterRuleRepository {
	return &BookContentFilterRuleRepository{db: db}
}

func (r *BookContentFilterRuleRepository) Create(ctx context.Context, m *model.BookContentFilterRule) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookContentFilterRuleRepository) Update(ctx context.Context, m *model.BookContentFilterRule) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookContentFilterRuleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookContentFilterRule{}, id).Error
}

func (r *BookContentFilterRuleRepository) GetByID(ctx context.Context, id uint64) (*model.BookContentFilterRule, error) {
	var m model.BookContentFilterRule
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookContentFilterRuleRepository) Page(ctx context.Context, req *dto.FilterRuleSearch) ([]model.BookContentFilterRule, int64, error) {
	var rows []model.BookContentFilterRule
	tx := r.db.WithContext(ctx).Model(&model.BookContentFilterRule{})
	if req.RuleName != "" {
		tx = tx.Where("rule_name LIKE ?", "%"+req.RuleName+"%")
	}
	if req.ApplyStage != "" {
		tx = tx.Where("apply_stage = ?", req.ApplyStage)
	}
	if req.Category != "" {
		tx = tx.Where("category = ?", req.Category)
	}
	if req.Status != "" {
		tx = tx.Where("status = ?", req.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("id DESC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *BookContentFilterRuleRepository) ListByStage(ctx context.Context, stage model.FilterApplyStage) ([]model.BookContentFilterRule, error) {
	var rows []model.BookContentFilterRule
	err := r.db.WithContext(ctx).Model(&model.BookContentFilterRule{}).
		Where("apply_stage = ? AND status = ?", stage, model.StatusEnabled).
		Find(&rows).Error
	return rows, err
}
