package em

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"

	"git.happyxhw.cn/happyxhw/iself/pkg/util"
)

// RequestID 获取并设置 request id
func RequestID() echo.MiddlewareFunc {
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

func GetCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), echo.HeaderXRequestID, c.Get(echo.HeaderXRequestID))
}
