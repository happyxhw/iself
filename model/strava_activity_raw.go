package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// StravaActivityRaw 原始活动信息
type StravaActivityRaw struct {
	ID        int64                 `gorm:"column:id;primary_key" json:"id"`
	Data      string                `gorm:"column:data" json:"data"`
	CreatedAt time.Time             `gorm:"column:created_at" json:"created_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *StravaActivityRaw) TableName() string {
	return "strava_activity_raw"
}
