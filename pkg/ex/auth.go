package ex

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID       int64
	SourceID int64
	Source   string
	Email    string
}

// AuthRequired middleware
func AuthRequired() echo.MiddlewareFunc {
	// 2. Return middleware handler
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			sess, _ := session.Get("session", c)
			if sess == nil {
				return ErrAuth
			}
			id, _ := sess.Values["id"].(float64)
			sourceID, _ := sess.Values["source_id"].(float64)
			if id == 0 && sourceID == 0 {
				return ErrAuth
			}
			source, _ := sess.Values["source"].(string)
			email, _ := sess.Values["email"].(string)
			c.Set("user", User{ID: int64(id), SourceID: int64(sourceID), Source: source, Email: email})
			return next(c)
		}
	}
}

func GetUser(c echo.Context) User {
	u, _ := c.Get("user").(User)
	return u
}
