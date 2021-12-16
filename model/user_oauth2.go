package model

import (
	"time"
)

// UserOauth2 model
type UserOauth2 struct {
	ID        int64  `gorm:"column:id;" json:"id"`
	SourceID  int64  `gorm:"source_id" json:"source_id"`
	Source    string `gorm:"source" json:"source"`
	AvatarURL string `gorm:"avatar_url" json:"avatar_url"`

	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt int64     `gorm:"column:deleted_at;default:0" json:"deleted_at"`
}

// TableName 表名
func (*UserOauth2) TableName() string {
	return "user_oauth2"
}
