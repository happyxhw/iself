package repo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

const (
	mockKey = "mockKey"
)

func TestCacher_GetBytes(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	content, _ := json.Marshal(&mockToken)
	mock.Regexp().ExpectGet(mockKey).SetVal(string(content))

	cr := NewCacher(rdb)

	data, err := cr.GetBytes(context.TODO(), mockKey)

	require.NoError(t, err)
	require.NotNil(t, data)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_GetString(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	content, _ := json.Marshal(&mockToken)
	mock.Regexp().ExpectGet(mockKey).SetVal(string(content))

	cr := Cacher{rdb: rdb}

	data, err := cr.GetString(context.TODO(), mockKey)

	require.NoError(t, err)
	require.NotEmpty(t, data)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_Get(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	content, _ := json.Marshal(&mockToken)
	mock.Regexp().ExpectGet(mockKey).SetVal(string(content))

	cr := Cacher{rdb: rdb}

	var data oauth2.Token
	err := cr.GetObject(context.TODO(), mockKey, &data)

	require.NoError(t, err)
	require.Equal(t, data.AccessToken, mockToken.AccessToken)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_Set(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet(mockKey, mockToken.AccessToken, 0).SetVal("ok")

	cr := Cacher{rdb: rdb}

	err := cr.Set(context.TODO(), mockKey, mockToken.AccessToken, 0)

	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_SetObject(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSet(mockKey, `\[.*?\]`, 0).SetVal("ok")

	cr := Cacher{rdb: rdb}

	err := cr.SetObject(context.TODO(), mockKey, mockToken, 0)

	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_SetNX(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSetNX(mockKey, mockToken.AccessToken, 0).SetVal(true)

	cr := Cacher{rdb: rdb}

	r, err := cr.SetNX(context.TODO(), mockKey, mockToken.AccessToken, 0)

	require.NoError(t, err)
	require.Equal(t, r, true)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_SetNXObject(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.Regexp().ExpectSetNX(mockKey, `\[.*?\]`, 0).SetVal(false)

	cr := Cacher{rdb: rdb}

	r, err := cr.SetNXObject(context.TODO(), mockKey, mockToken, 0)

	require.NoError(t, err)
	require.Equal(t, r, false)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCacher_Del(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.Regexp().ExpectDel(mockKey).SetVal(1)

	cr := Cacher{rdb: rdb}

	r, err := cr.Del(context.TODO(), mockKey)

	require.NoError(t, err)
	require.Equal(t, r, int64(1))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
