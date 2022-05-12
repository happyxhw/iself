package em

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID     int64
	Source string
	Email  string
}

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
			id = 19830262
			source, _ := sess.Values["source"].(string)
			email, _ := sess.Values["email"].(string)
			c.Set("user", User{ID: int64(id), Source: source, Email: email})
			return next(c)
		}
	}
}

func GetUser(c echo.Context) User {
	u, _ := c.Get("user").(User)
	return u
}
