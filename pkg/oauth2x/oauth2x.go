package oauth2x

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"git.happyxhw.cn/happyxhw/iself/model"
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
