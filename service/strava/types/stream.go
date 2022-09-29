package types

import (
	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/strava"
)

type StreamSet struct {
	*strava.StreamSet

	ID int64
}

func NewStreamSet(m *model.StravaActivityStream) *StreamSet {
	return &StreamSet{
		ID: m.ID,
		StreamSet: &strava.StreamSet{
			Time:           m.TimeStream,
			Distance:       m.DistanceStream,
			Latlng:         m.LatlngStream,
			Altitude:       m.AltitudeStream,
			VelocitySmooth: m.VelocitySmoothStream,
			Heartrate:      m.HeartrateStream,
			Cadence:        m.CadenceStream,
			Watts:          m.WattsStream,
			Temp:           m.TempStream,
			Moving:         m.MovingStream,
			GradeSmooth:    m.GradeSmoothStream,
		},
	}
}
