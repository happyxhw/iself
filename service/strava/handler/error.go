package handler

import (
	"net/http"

	"github.com/happyxhw/iself/pkg/ex"
)

var (
	ErrStravaAPI   = ex.NewError(http.StatusServiceUnavailable, 240001, "strava api error")
	ErrStravaToken = ex.NewError(http.StatusServiceUnavailable, 240002, "strava oauth2 token error")
)
