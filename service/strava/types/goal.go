package types

type GoalReq struct {
	SourceID int64
	Type     string  `query:"type" validate:"activity"`
	Field    string  `query:"field" validate:"stats_field"`
	Freq     string  `query:"freq" validate:"oneof=week month year"`
	Value    float64 `query:"freq" validate:"gte=1"`
}

type QueryGoalReq struct {
	SourceID int64
	Type     string `query:"type" validate:"activity"`
	Field    string `query:"field" validate:"stats_field"`
}

type QueryGoalResp struct {
	Type  string  `json:"type"`
	Field string  `json:"field"`
	Freq  string  `json:"freq"`
	Value float64 `json:"value"`
}
