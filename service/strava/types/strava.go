package types

const (
	All         = "all"
	Run         = "run"
	Ride        = "ride"
	VirtualRide = "virtualride"
)

const (
	Week  = "week"
	Month = "month"
	Year  = "year"
)

// ActivityListReq request for activity list
type ActivityListReq struct {
	Limit  int   `query:"limit"`
	After  int64 `query:"after"`
	Before int64 `query:"before"`

	Type string `query:"type" validate:"activity"`
}

// ActivityListPageReq request for activity list
type ActivityListPageReq struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Type     string `query:"type" validate:"activity"`
}

type StatsNowReq struct {
	Field  string `query:"field" validate:"stats_field"`
	Type   string `query:"type" validate:"activity"`
	Method string `query:"method" validate:"stats_method"`
}

type DateChartReq struct {
	Field  string `query:"field" validate:"stats_field"`
	Type   string `query:"type" validate:"activity"`
	Method string `query:"method" validate:"stats_method"`
	Freq   string `query:"freq" validate:"oneof=week month year"`
	Size   int    `query:"size"`
}

// Stats sum() | average | max | min group by week, month, year stats
type Stats struct {
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
