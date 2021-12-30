package em

import (
	"github.com/go-playground/validator"
)

// AlpineSki, BackcountrySki, Canoeing, Crossfit, EBikeRide, Elliptical, Golf, Handcycle, Hike, IceSkate,
// InlineSkate, Kayaking, Kitesurf, NordicSki, Ride, RockClimbing, RollerSki, Rowing, Run,
// Sail, Skateboard, Snowboard, Snowshoe, Soccer, StairStepper, StandUpPaddling, Surfing, Swim,
// Velomobile, VirtualRide, VirtualRun, Walk, WeightTraining, Wheelchair, Windsurf, Workout, Yoga

var typeMap = map[string]bool{
	"run":         true,
	"ride":        true,
	"virtualride": true,
	"all":         true,
}

var fieldMap = map[string]bool{
	"distance":    true,
	"calories":    true,
	"moving_time": true,
}

var freqMap = map[string]bool{
	"week":  true,
	"month": true,
	"year":  true,
	"all":   true,
}

var methodMap = map[string]bool{
	"sum": true, "max": true, "min": true, "avg": true,
}

// ActivityType validate activity type
func ActivityType(fl validator.FieldLevel) bool {
	activityType := fl.Field().String()
	return typeMap[activityType]
}

// StatsField validate stats field
func StatsField(fl validator.FieldLevel) bool {
	return fieldMap[fl.Field().String()]
}

// StatsFreq validate stats freq
func StatsFreq(fl validator.FieldLevel) bool {
	return freqMap[fl.Field().String()]
}

// StatsMethod validate stats method
func StatsMethod(fl validator.FieldLevel) bool {
	return methodMap[fl.Field().String()]
}

type (
	CustomValidator struct {
		Validator *validator.Validate
	}
)

func NewValidator() *CustomValidator {
	v := CustomValidator{
		Validator: validator.New(),
	}
	_ = v.Validator.RegisterValidation("activity", ActivityType)
	_ = v.Validator.RegisterValidation("stats_field", StatsField)
	_ = v.Validator.RegisterValidation("stats_freq", StatsFreq)
	_ = v.Validator.RegisterValidation("stats_method", StatsMethod)

	return &v
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
