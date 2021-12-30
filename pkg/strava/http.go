package strava

import (
	"context"
	"io"
	"net/http"

	"go.uber.org/zap"

	"git.happyxhw.cn/happyxhw/iself/pkg/log"
)

// do make request to strava
func do(ctx context.Context, url, method string, body io.Reader, client *http.Client) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		log.Error("new request", zap.Error(err))
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("do request", zap.Error(err))
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)
}
