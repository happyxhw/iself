package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	activityAPI = "/activities"

	streamSet = "time,distance,latlng,altitude,velocity_smooth,heartrate,cadence,watts,temp,moving,grade_smooth"
	streamAPI = "/activities/%d/streams?key_by_type=true&keys=%s"
)

type Activity service

// Activity get activity info
func (s *Activity) Activity(ctx context.Context, id int64) (*DetailedActivity, []byte, error) {
	url := fmt.Sprintf("%s%s/%d", s.client.BaseURL, activityAPI, id)
	body, err := do(ctx, url, http.MethodGet, http.NoBody, s.client.httpClient)
	if err != nil {
		return nil, nil, err
	}
	var resp DetailedActivity
	err = json.Unmarshal(body, &resp)
	return &resp, body, err
}

// ActivityStream get activity stream
func (s *Activity) ActivityStream(ctx context.Context, id int64) (*StreamSet, error) {
	api := fmt.Sprintf(streamAPI, id, streamSet)
	url := fmt.Sprintf("%s%s", s.client.BaseURL, api)
	body, err := do(ctx, url, http.MethodGet, http.NoBody, s.client.httpClient)
	if err != nil {
		return nil, err
	}
	var resp StreamSet
	err = json.Unmarshal(body, &resp)
	return &resp, err
}
