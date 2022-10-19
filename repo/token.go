package repo

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

const (
	refreshTimeout = time.Second * 60
)

// TokenRepo srv
type TokenRepo struct {
	cacher *Cacher
}

func NewTokenRepo(cacher *Cacher) *TokenRepo {
	return &TokenRepo{
		cacher: cacher,
	}
}

func (tr *TokenRepo) SaveToken(ctx context.Context, token *oauth2.Token, source string, sourceID int64) error {
	key := fmt.Sprintf("oauth2:%s:%d", source, sourceID)

	return tr.cacher.SetObject(ctx, key, token, 0)
}

type Refresher func(context.Context, *oauth2.Token) (*oauth2.Token, error)

func (tr *TokenRepo) GetToken(ctx context.Context, source string, sourceID int64, fn Refresher) (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, refreshTimeout)
	defer cancel()

	var token oauth2.Token
	key := fmt.Sprintf("oauth2:%s:%d", source, sourceID)
	err := tr.cacher.GetObject(ctx, key, &token)
	if err != nil {
		return nil, err
	}
	refreshedToken, err := fn(ctx, &token)
	if err != nil {
		return nil, err
	}
	if refreshedToken.AccessToken != token.AccessToken {
		if err := tr.SaveToken(ctx, refreshedToken, source, sourceID); err != nil {
			return nil, err
		}
	}

	return refreshedToken, nil
}
