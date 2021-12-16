package em

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// AuthRequired middleware
func AuthRequired() echo.MiddlewareFunc {
	// 2. Return middleware handler
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			sess, _ := session.Get("session", c)
			id, _ := sess.Values["id"].(float64)
			if int64(id) == 0 {
				return ErrAuth
			}
			source, _ := sess.Values["source"].(string)
			c.Set("id", int64(id))
			c.Set("source", source)
			return next(c)
		}
	}
}
