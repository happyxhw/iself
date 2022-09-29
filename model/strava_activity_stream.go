package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"

	"git.happyxhw.cn/happyxhw/iself/pkg/strava"
)

// StravaActivityStream stream 信息
type StravaActivityStream struct {
	ID             int64                 `gorm:"column:id;primary_key" json:"id"`
	Time           []byte                `gorm:"column:time" json:"time"`
	Distance       []byte                `gorm:"column:distance" json:"distance"`
	Latlng         []byte                `gorm:"column:latlng" json:"latlng"`
	Altitude       []byte                `gorm:"column:altitude" json:"altitude"`
	VelocitySmooth []byte                `gorm:"column:velocity_smooth" json:"velocity_smooth"`
	Heartrate      []byte                `gorm:"column:heartrate" json:"heartrate"`
	Cadence        []byte                `gorm:"column:cadence" json:"cadence"`
	Watts          []byte                `gorm:"column:watts" json:"watts"`
	Temp           []byte                `gorm:"column:temp" json:"temp"`
	Moving         []byte                `gorm:"column:moving" json:"moving"`
	GradeSmooth    []byte                `gorm:"column:grade_smooth" json:"grade_smooth"`
	CreatedAt      time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time             `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedAt      soft_delete.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	TimeStream           *strava.TimeStream           `gorm:"-"`
	DistanceStream       *strava.DistanceStream       `gorm:"-"`
	LatlngStream         *strava.LatLngStream         `gorm:"-"`
	AltitudeStream       *strava.AltitudeStream       `gorm:"-"`
	VelocitySmoothStream *strava.SmoothVelocityStream `gorm:"-"`
	HeartrateStream      *strava.HeartrateStream      `gorm:"-"`
	CadenceStream        *strava.CadenceStream        `gorm:"-"`
	WattsStream          *strava.PowerStream          `gorm:"-"`
	TempStream           *strava.TemperatureStream    `gorm:"-"`
	MovingStream         *strava.MovingStream         `gorm:"-"`
	GradeSmoothStream    *strava.SmoothGradeStream    `gorm:"-"`
}

func (s *StravaActivityStream) TableName() string {
	return "strava_activity_stream"
}

//nolint:gocyclo,funlen
func (s *StravaActivityStream) BeforeCreate(tx *gorm.DB) error {
	if s.TimeStream != nil {
		data, err := json.Marshal(s.TimeStream)
		if err != nil {
			return err
		}
		s.Time = data
	}
	if s.DistanceStream != nil {
		data, err := json.Marshal(s.DistanceStream)
		if err != nil {
			return err
		}
		s.Distance = data
	}

	if s.LatlngStream != nil {
		data, err := json.Marshal(s.LatlngStream)
		if err != nil {
			return err
		}
		s.Latlng = data
	}
	if s.AltitudeStream != nil {
		data, err := json.Marshal(s.AltitudeStream)
		if err != nil {
			return err
		}
		s.Altitude = data
	}
	if s.VelocitySmoothStream != nil {
		data, err := json.Marshal(s.VelocitySmoothStream)
		if err != nil {
			return err
		}
		s.VelocitySmooth = data
	}
	if s.HeartrateStream != nil {
		data, err := json.Marshal(s.HeartrateStream)
		if err != nil {
			return err
		}
		s.Heartrate = data
	}
	if s.CadenceStream != nil {
		data, err := json.Marshal(s.CadenceStream)
		if err != nil {
			return err
		}
		s.Cadence = data
	}
	if s.WattsStream != nil {
		data, err := json.Marshal(s.WattsStream)
		if err != nil {
			return err
		}
		s.Watts = data
	}
	if s.TempStream != nil {
		data, err := json.Marshal(s.TempStream)
		if err != nil {
			return err
		}
		s.Temp = data
	}
	if s.MovingStream != nil {
		data, err := json.Marshal(s.MovingStream)
		if err != nil {
			return err
		}
		s.Moving = data
	}
	if s.GradeSmoothStream != nil {
		data, err := json.Marshal(s.GradeSmoothStream)
		if err != nil {
			return err
		}
		s.GradeSmooth = data
	}

	return nil
}

//nolint:gocyclo
func (s *StravaActivityStream) AfterFind(tx *gorm.DB) error {
	if len(s.Time) > 0 {
		s.TimeStream = new(strava.TimeStream)
		if err := json.Unmarshal(s.Time, s.TimeStream); err != nil {
			return err
		}
	}
	if len(s.Distance) > 0 {
		s.DistanceStream = new(strava.DistanceStream)
		if err := json.Unmarshal(s.Distance, s.DistanceStream); err != nil {
			return err
		}
	}
	if len(s.Latlng) > 0 {
		s.LatlngStream = new(strava.LatLngStream)
		if err := json.Unmarshal(s.Latlng, s.LatlngStream); err != nil {
			return err
		}
	}
	if len(s.Altitude) > 0 {
		s.AltitudeStream = new(strava.AltitudeStream)
		if err := json.Unmarshal(s.Altitude, s.AltitudeStream); err != nil {
			return err
		}
	}
	if len(s.VelocitySmooth) > 0 {
		s.VelocitySmoothStream = new(strava.SmoothVelocityStream)
		if err := json.Unmarshal(s.VelocitySmooth, s.VelocitySmoothStream); err != nil {
			return err
		}
	}
	if len(s.Heartrate) > 0 {
		s.HeartrateStream = new(strava.HeartrateStream)
		if err := json.Unmarshal(s.Heartrate, s.HeartrateStream); err != nil {
			return err
		}
	}
	if len(s.Cadence) > 0 {
		s.CadenceStream = new(strava.CadenceStream)
		if err := json.Unmarshal(s.Cadence, s.CadenceStream); err != nil {
			return err
		}
	}
	if len(s.Watts) > 0 {
		s.WattsStream = new(strava.PowerStream)
		if err := json.Unmarshal(s.Watts, s.WattsStream); err != nil {
			return err
		}
	}
	if len(s.Temp) > 0 {
		s.TempStream = new(strava.TemperatureStream)
		if err := json.Unmarshal(s.Temp, s.TempStream); err != nil {
			return err
		}
	}
	if len(s.Moving) > 0 {
		s.MovingStream = new(strava.MovingStream)
		if err := json.Unmarshal(s.Moving, s.MovingStream); err != nil {
			return err
		}
	}
	if len(s.GradeSmooth) > 0 {
		s.GradeSmoothStream = new(strava.SmoothGradeStream)
		if err := json.Unmarshal(s.GradeSmooth, s.GradeSmoothStream); err != nil {
			return err
		}
	}

	return nil
}
