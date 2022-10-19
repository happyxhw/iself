package types

import "github.com/happyxhw/iself/model"

type Goal struct {
	ID    int64   `json:"id"`
	Type  string  `json:"type"`
	Field string  `json:"field"`
	Freq  string  `json:"freq"`
	Value float64 `json:"value"`

	AthleteID int64 `json:"-"`
}

func NewGoal(m *model.StravaGoal) *Goal {
	return &Goal{
		ID:    m.ID,
		Type:  m.Type,
		Field: m.Field,
		Freq:  m.Freq,
		Value: m.Value,
	}
}

func NewGoals(ms []*model.StravaGoal) []*Goal {
	list := make([]*Goal, 0, len(ms))
	for _, item := range ms {
		list = append(list, NewGoal(item))
	}
	return list
}

type CreateGoalReq struct {
	Type  string  `query:"type" validate:"activity"`
	Field string  `query:"field" validate:"stats_field"`
	Freq  string  `query:"freq" validate:"oneof=week month year"`
	Value float64 `query:"freq" validate:"gte=1"`
}

type UpdateGoalReq struct {
	ID    int64   `param:"id"`
	Value float64 `query:"freq" validate:"gte=1"`
}

type QueryGoalReq struct {
	Type  string `query:"type" validate:"activity"`
	Field string `query:"field" validate:"stats_field"`
}
