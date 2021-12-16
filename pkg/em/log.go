package em

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Logger for api access
func Logger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			req := c.Request()
			resp := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = resp.Header().Get(echo.HeaderXRequestID)
			}
			latency := time.Since(start).Milliseconds()

			status := resp.Status
			reqSize, _ := strconv.ParseInt(req.Header.Get(echo.HeaderContentLength), 10, 64)
			respSize := resp.Size

			var errCode int
			var errMsg string
			if he, ok := err.(*Error); ok {
				errCode = he.Code
				errMsg = he.Error()
			} else {
				errMsg = err.Error()
			}

			fields := []zap.Field{
				zap.Int("code", status),
				zap.Int("err_code", errCode),
				zap.String("method", req.Method),
				zap.Int64("latency", latency),
				zap.String("ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("refer", req.Referer()),
				zap.String("path", c.Path()),
				zap.String("uri", req.RequestURI),
				zap.Int64("req_size", reqSize),
				zap.Int64("resp_size", respSize),
				zap.String("err", errMsg),
				zap.String("request_id", id),
				zap.String("ua", req.UserAgent()),
			}

			switch {
			case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
				logger.Warn("[API]", fields...)
			case status >= http.StatusInternalServerError:
				logger.Error("[API]", fields...)
			default:
				logger.Info("[API]", fields...)
			}

			return nil
		}
	}
}
