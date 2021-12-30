package strava

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	athleteAPI = "/athlete"
)

type Athlete service

// Athlete get athlete basic info
func (s *Athlete) Athlete(ctx context.Context) (*SummaryAthlete, error) {
	url := s.client.BaseURL + athleteAPI
	body, err := do(ctx, url, http.MethodGet, http.NoBody, s.client.httpClient)
	if err != nil {
		return nil, err
	}
	var resp SummaryAthlete
	err = json.Unmarshal(body, &resp)
	return &resp, err
}
