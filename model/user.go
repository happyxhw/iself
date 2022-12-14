package model

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"github.com/happyxhw/pkg/query"
)

// User model
type User struct {
	ID        int64  `gorm:"column:id;" json:"id"`
	Name      string `gorm:"column:name;" json:"name"`
	Email     string `gorm:"column:email;" json:"email"`
	Password  string `gorm:"column:password;" json:"password"`
	AvatarURL string `gorm:"avatar_url" json:"avatar_url"`
	Role      int    `gorm:"role" json:"role"`
	Source    string `gorm:"source" json:"source"`
	SourceID  int64  `gorm:"source_id" json:"source_id"`
	Status    int    `gorm:"status" json:"status"`

	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;default:0" json:"deleted_at"`
}

// TableName 表名
func (*User) TableName() string {
	return "user"
}

type UserStatus int

const (
	WaitActiveStatus UserStatus = iota
	ActivatedStatus
)

type UserParam struct {
	query.Param `gorm:"-"`

	Name      *string `gorm:"column:name;" json:"name"`
	Password  *string `gorm:"column:password;" json:"password"`
	AvatarURL *string `gorm:"avatar_url" json:"avatar_url"`
	Role      *int    `gorm:"role" json:"role"`
	Status    *int    `gorm:"status" json:"status"`

	UpdatedAt *time.Time            `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;default:0" json:"deleted_at"`
}
