package ex

import (
	"github.com/happyxhw/pkg/util"
	"github.com/labstack/echo/v4"
)

// RequestID 获取并设置 request id
func SetRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = util.NanoID(16)
			}
			c.Set(echo.HeaderXRequestID, requestID)

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = util.NanoID(16)
			}
			c.Set(echo.HeaderXRequestID, rid)
			res.Header().Set(echo.HeaderXRequestID, rid)
			return next(c)
		}
	}
}
