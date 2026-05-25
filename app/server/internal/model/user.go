package model

type User struct {
	BaseModel
	Username string `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(128);not null" json:"-"`
	Nickname string `gorm:"type:varchar(64)" json:"nickname"`
	Email    string `gorm:"type:varchar(128)" json:"email"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	Avatar   string `gorm:"type:varchar(256)" json:"avatar"`
	Status   int    `gorm:"type:tinyint;default:1" json:"status"`
}

func (User) TableName() string {
	return "users"
}