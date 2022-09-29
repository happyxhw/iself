package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// StravaGoal model
type StravaGoal struct {
	ID        int64   `gorm:"column:id;" json:"id"`
	AthleteID int64   `gorm:"column:athlete_id" json:"athlete_id"`
	Type      string  `gorm:"column:type" json:"type"`
	Field     string  `gorm:"column:field" json:"field"`
	Freq      string  `gorm:"column:freq" json:"freq"`
	Value     float64 `gorm:"column:value" json:"value"`

	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;default:0" json:"deleted_at"`
}

// TableName 表名
func (*StravaGoal) TableName() string {
	return "strava_goal"
}

type StravaGoalParam struct {
	Value *float64 `gorm:"column:value" json:"value"`

	UpdatedAt *time.Time            `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;default:0" json:"deleted_at"`
}
