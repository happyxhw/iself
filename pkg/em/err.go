package em

import (
	"github.com/labstack/echo/v4"
)

var _ error = (*Error)(nil)

// Error 封装接口错误信息
type Error struct {
	Status       int    `json:"-"`               // Status http 状态码
	Code         int    `json:"code"`            // Code 错误状态码
	Message      string `json:"message"`         // Message 错误信息
	ErrorMessage string `json:"error,omitempty"` // ErrorMessage 只有在非 gin.ReleaseMode 模式下才会在接口中展示，但会在日志中展示
}

// NewError return a error wrap
func NewError(status, code int, message string) *Error {
	e := &Error{
		Status:  status,
		Code:    code,
		Message: message,
	}
	return e
}

// Error 返回错误信息
func (e *Error) Error() string {
	if e.ErrorMessage != "" {
		return e.ErrorMessage
	}
	return e.Message
}

// Wrap 包装 err
func (e *Error) Wrap(err error) *Error {
	if err != nil {
		e.ErrorMessage = err.Error()
	}
	return e
}

// ErrHandler handler
func ErrHandler(e *echo.Echo) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		switch v := (err).(type) {
		case *Error:
			if !e.Debug {
				_ = c.JSON(v.Status, Error{Code: v.Code, Message: v.Message})
				return
			}
			_ = c.JSON(v.Status, v)
			return
		case *echo.HTTPError:
			_ = c.JSON(v.Code, v)
			return
		}

		_ = c.JSON(c.Response().Status, err)
	}
}
