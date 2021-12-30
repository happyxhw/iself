package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// Goal model
type Goal struct {
	ID        int64   `gorm:"column:id;" json:"id"`
	AthleteID int64   `gorm:"column:athlete_id" json:"athlete_id"`
	Type      string  `gorm:"column:type" json:"type"`
	Field     string  `gorm:"column:field" json:"field"`
	Freq      string  `gorm:"column:freq" json:"freq"`
	Value     float64 `gorm:"column:value" json:"value"`

	CreatedAt time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 表名
func (*Goal) TableName() string {
	return "goal"
}
