package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

var mockToken = oauth2.Token{
	AccessToken:  "at",
	TokenType:    "code",
	RefreshToken: "rk",
}

func mockRefresher(_ context.Context, _ *oauth2.Token) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "new_at",
		TokenType:    "code",
		RefreshToken: "rk",
	}, nil
}

func TestTokenRepo_SaveToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	key := fmt.Sprintf("oauth2:%s:%d", mockUser.Source, mockUser.ID)
	mock.Regexp().ExpectSet(key, `\[.*?\]`, 0).SetVal("ok")

	tr := NewTokenRepo(NewCacher(rdb))

	err := tr.SaveToken(context.TODO(), nil, mockUser.Source, mockUser.ID)

	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestTokenRepo_GetToken(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	key := fmt.Sprintf("oauth2:%s:%d", mockUser.Source, mockUser.ID)
	content, _ := json.Marshal(&mockToken)
	mock.Regexp().ExpectGet(key).SetVal(string(content))
	mock.Regexp().ExpectSet(key, `\[.*?\]`, 0).SetVal("ok")

	tr := NewTokenRepo(NewCacher(rdb))

	token, err := tr.GetToken(context.TODO(), mockUser.Source, mockUser.ID, mockRefresher)

	require.NoError(t, err)
	require.NotNil(t, token)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
