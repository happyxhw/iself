package types

// ActivityProgressStats 相关数据聚合统计, 如: 本周跑步里程, 本月消耗卡路里等
// Stats sum() | average | max | min group by week, month, year stats
type ActivityProgressStats struct {
	Unit string `json:"unit"`
	Type string `json:"type"`

	All string `json:"all"`

	Week        string `json:"week"`
	WeekGoal    string `json:"week_goal"`
	WeekProcess string `json:"week_process"`

	Month        string `json:"month"`
	MonthGoal    string `json:"month_goal"`
	MonthProcess string `json:"month_process"`

	Year        string `json:"year"`
	YearGoal    string `json:"year_goal"`
	YearProcess string `json:"year_process"`
}

type ActivityAggStats struct {
	Value    []float64 `json:"value"`
	Time     []string  `json:"time"`
	Max      float64   `json:"max"`
	Min      float64   `json:"min"`
	Avg      float64   `json:"avg"`
	MaxIndex int       `json:"max_index"`
	MinIndex int       `json:"min_index"`
}

type ProgressStatsReq struct {
	Field  string `query:"field" validate:"stats_field"`
	Type   string `query:"type" validate:"activity"`
	Method string `query:"method" validate:"stats_method"`
}

type AggStatsReq struct {
	Field  string `query:"field" validate:"stats_field"`
	Type   string `query:"type" validate:"activity"`
	Method string `query:"method" validate:"stats_method"`
	Freq   string `query:"freq" validate:"oneof=week month year"`
	Size   int    `query:"size"`
}
