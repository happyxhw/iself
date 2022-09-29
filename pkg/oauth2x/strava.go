package oauth2x

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/strava"
)

type Strava struct {
	conf *oauth2.Config
}

func NewStrava(conf *oauth2.Config) *Strava {
	return &Strava{
		conf: conf,
	}
}

func (s *Strava) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.conf.Exchange(ctx, code)
	return token, err
}

func (s *Strava) GetUser(ctx context.Context, token *oauth2.Token) (*model.User, error) {
	cli := s.conf.Client(ctx, token)
	stravaCli := strava.NewClient(cli)
	athlete, err := stravaCli.Athlete.Athlete(ctx)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Name:      athlete.Username,
		Source:    StravaSource,
		SourceID:  athlete.Id,
		AvatarURL: athlete.ProfileMedium,
	}
	return &user, nil
}

func (s *Strava) Refresh(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := s.conf.TokenSource(ctx, token).Token()
	return newToken, err
}

func (s *Strava) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return s.conf.Client(ctx, token)
}
