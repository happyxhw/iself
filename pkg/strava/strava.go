package strava

import (
	"net/http"
)

const (
	BaseURL = "https://www.strava.com/api/v3"
)

type Client struct {
	httpClient *http.Client

	BaseURL string

	Athlete  *Athlete
	Activity *Activity

	common service
}

func NewClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient: httpClient,
		BaseURL:    BaseURL,
	}
	c.common.client = c
	c.Athlete = (*Athlete)(&c.common)
	c.Activity = (*Activity)(&c.common)

	return c
}

type service struct {
	client *Client
}
