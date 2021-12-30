package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// ActivityDetail 活动详细信息
type ActivityDetail struct {
	ID                 int64                 `gorm:"column:id;primary_key" json:"id"`
	AthleteID          int64                 `gorm:"column:athlete_id;NOT NULL" json:"athlete_id"`
	Name               string                `gorm:"column:name;NOT NULL" json:"name"`
	Type               string                `gorm:"column:type;NOT NULL" json:"type"`
	Distance           float64               `gorm:"column:distance;default:0.0;NOT NULL" json:"distance"`
	MovingTime         int                   `gorm:"column:moving_time;default:0;NOT NULL" json:"moving_time"`
	ElapsedTime        int                   `gorm:"column:elapsed_time;default:0;NOT NULL" json:"elapsed_time"`
	TotalElevationGain float64               `gorm:"column:total_elevation_gain;default:0.0;NOT NULL" json:"total_elevation_gain"`
	StartDateLocal     time.Time             `gorm:"column:start_date_local;NOT NULL" json:"start_date_local"`
	Polyline           string                `gorm:"column:polyline;NOT NULL" json:"polyline"`
	SummaryPolyline    string                `gorm:"column:summary_polyline;NOT NULL" json:"summary_polyline"`
	AverageSpeed       float64               `gorm:"column:average_speed;default:0.0;NOT NULL" json:"average_speed"`
	MaxSpeed           float64               `gorm:"column:max_speed;default:0.0;NOT NULL" json:"max_speed"`
	AverageHeartrate   float64               `gorm:"column:average_heartrate;default:0.0;NOT NULL" json:"average_heartrate"`
	MaxHeartrate       float64               `gorm:"column:max_heartrate;default:0.0;NOT NULL" json:"max_heartrate"`
	ElevHigh           float64               `gorm:"column:elev_high;default:0.0;NOT NULL" json:"elev_high"`
	ElevLow            float64               `gorm:"column:elev_low;default:0.0;NOT NULL" json:"elev_low"`
	Calories           float64               `gorm:"column:calories;default:0.0;NOT NULL" json:"calories"`
	SplitsMetric       string                `gorm:"column:splits_metric" json:"splits_metric"`
	BestEfforts        string                `gorm:"column:best_efforts" json:"best_efforts"`
	DeviceName         string                `gorm:"column:device_name" json:"device_name"`
	CreatedAt          *time.Time            `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt          *time.Time            `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt          soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,,omitempty"`
}

func (m *ActivityDetail) TableName() string {
	return "activity_detail"
}

// ActivityRaw 原始活动信息
type ActivityRaw struct {
	ID        int64                 `gorm:"column:id;primary_key" json:"id"`
	Data      string                `gorm:"column:data" json:"data"`
	CreatedAt time.Time             `gorm:"column:created_at" json:"created_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *ActivityRaw) TableName() string {
	return "activity_raw"
}

// ActivityStream stream 信息
type ActivityStream struct {
	ID             int64                 `gorm:"column:id;primary_key" json:"id"`
	Time           string                `gorm:"column:time" json:"time"`
	Distance       string                `gorm:"column:distance" json:"distance"`
	Latlng         string                `gorm:"column:latlng" json:"latlng"`
	Altitude       string                `gorm:"column:altitude" json:"altitude"`
	VelocitySmooth string                `gorm:"column:velocity_smooth" json:"velocity_smooth"`
	Heartrate      string                `gorm:"column:heartrate" json:"heartrate"`
	Cadence        string                `gorm:"column:cadence" json:"cadence"`
	Watts          string                `gorm:"column:watts" json:"watts"`
	Temp           string                `gorm:"column:temp" json:"temp"`
	Moving         string                `gorm:"column:moving" json:"moving"`
	GradeSmooth    string                `gorm:"column:grade_smooth" json:"grade_smooth"`
	CreatedAt      time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time             `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt      soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *ActivityStream) TableName() string {
	return "activity_stream"
}
