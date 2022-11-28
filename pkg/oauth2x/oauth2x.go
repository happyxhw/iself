package oauth2x

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/happyxhw/iself/model"
)

const (
	StravaSource = "strava"
	GithubSource = "github"
)

type Oauth2x interface {
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (*model.User, error)
	Refresh(context.Context, *oauth2.Token) (*oauth2.Token, error)
	Client(ctx context.Context, token *oauth2.Token) *http.Client
}

type ClientConfig struct {
	Name         string
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Scopes       []string
}

var provider map[string]Oauth2x

func InitProvider(cfg []*ClientConfig) {
	provider = make(map[string]Oauth2x, len(cfg))
	for _, c := range cfg {
		if c.Name == StravaSource {
			conf := oauth2.Config{
				ClientID:     c.ClientID,
				ClientSecret: c.ClientSecret,
				Endpoint:     endpoints.Strava,
				Scopes:       c.Scopes,
			}
			provider[c.Name] = NewStrava(&conf)
		}
	}
}

func Provider() map[string]Oauth2x {
	return provider
}
