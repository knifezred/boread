package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

// BookChapterRuleRepository 章节识别规则数据访问层
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

// BookChapterRuleRelRepository 书籍章节规则关联数据访问层
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
