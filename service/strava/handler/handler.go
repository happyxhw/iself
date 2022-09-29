package handler

const (
	notExistsLabel = "--"
)

var unitMap = map[string]string{
	"distance":    "km",
	"moving_time": "s",
	"calories":    "Cal",
}

var fractionMap = map[string]float64{
	"distance": 1000,
}

const (
	limitWeek  = 12
	limitMonth = 12
	limitYear  = 12
)

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
