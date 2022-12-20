package model

import (
	"encoding/json"
	"time"

	"github.com/happyxhw/pkg/query"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"

	"github.com/happyxhw/iself/pkg/strava"
)

// StravaActivityDetail 活动详细信息
type StravaActivityDetail struct {
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
	SplitsMetric       []byte                `gorm:"column:splits_metric" json:"splits_metric"`
	BestEfforts        []byte                `gorm:"column:best_efforts" json:"best_efforts"`
	DeviceName         string                `gorm:"column:device_name" json:"device_name"`
	CreatedAt          time.Time             `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt          time.Time             `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt          soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,,omitempty"`

	SplitsMetricJSON []*strava.Split                 `gorm:"-"`
	BestEffortsJSON  []*strava.DetailedSegmentEffort `gorm:"-"`
}

func (m *StravaActivityDetail) TableName() string {
	return "strava_activity_detail"
}

type StravaActivityParam struct {
	query.Param `gorm:"-"`

	Type      *string               `gorm:"column:type;NOT NULL" json:"type"`
	UpdatedAt *time.Time            `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,,omitempty"`
}

func (m *StravaActivityDetail) BeforeCreate(tx *gorm.DB) error {
	if len(m.SplitsMetricJSON) > 0 {
		data, err := json.Marshal(m.SplitsMetricJSON)
		if err != nil {
			return err
		}
		m.SplitsMetric = data
	}
	if len(m.BestEffortsJSON) > 0 {
		data, err := json.Marshal(m.BestEffortsJSON)
		if err != nil {
			return err
		}
		m.BestEfforts = data
	}

	return nil
}

func (m *StravaActivityDetail) AfterFind(tx *gorm.DB) error {
	if len(m.SplitsMetric) > 0 {
		if err := json.Unmarshal(m.SplitsMetric, &m.SplitsMetricJSON); err != nil {
			return err
		}
	}
	if len(m.BestEfforts) > 0 {
		if err := json.Unmarshal(m.BestEfforts, &m.BestEffortsJSON); err != nil {
			return err
		}
	}

	return nil
}
