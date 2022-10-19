package types

import (
	"github.com/jinzhu/copier"

	"github.com/happyxhw/iself/model"
	"github.com/happyxhw/iself/pkg/query"
	"github.com/happyxhw/iself/pkg/strava"
)

type Activity struct {
	DetailedActivity *DetailedActivity `json:"detailed_activity"`
	StreamSet        *StreamSet        `json:"stream_set"`
}

type DetailedActivity struct {
	*strava.DetailedActivity
}

func NewDetailedActivity(m *model.StravaActivityDetail) *DetailedActivity {
	var a DetailedActivity
	_ = copier.Copy(&a, m)
	if m.Polyline != "" {
		a.Map = &strava.PolylineMap{
			Polyline:        m.Polyline,
			SummaryPolyline: m.SummaryPolyline,
		}
	}
	return &a
}

func NewDetailedActivityList(s []*model.StravaActivityDetail) []*DetailedActivity {
	list := make([]*DetailedActivity, 0, len(s))
	for _, item := range s {
		list = append(list, NewDetailedActivity(item))
	}

	return list
}

type ActivityQueryParam struct {
	query.Param
	ActivityType *string `query:"type"`
}

type ActivityQueryResult struct {
	PageResult *query.PagingResult `json:"page"`
	Data       []*DetailedActivity `json:"data"`
}
