package model

import (
	"time"
)

// PushEvent 推送事件
type PushEvent struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`
	AspectType string    `gorm:"column:aspect_type" json:"aspect_type"`
	EventTime  int64     `gorm:"column:event_time" json:"event_time"`
	ObjectID   int64     `gorm:"column:object_id" json:"object_id"`
	ObjectType string    `gorm:"column:object_type" json:"object_type"`
	OwnerID    int64     `gorm:"column:owner_id" json:"owner_id"`
	Updates    string    `gorm:"column:updates" json:"updates"`
	Status     int       `gorm:"column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PushEvent) TableName() string {
	return "push_event"
}
