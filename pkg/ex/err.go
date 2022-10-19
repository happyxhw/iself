package ex

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap/zapcore"
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

// Msg wrap msg
func (e *Error) Msg(msg string) *Error {
	newErr := Error{
		Status:       e.Status,
		Code:         e.Code,
		Message:      e.Message,
		ErrorMessage: e.ErrorMessage,
	}
	if msg != "" {
		if newErr.Message == "" {
			newErr.Message = msg
		} else {
			newErr.Message = fmt.Sprintf("%s; %s", newErr.Message, msg)
		}
	}
	return &newErr
}

// Wrap 包装 err
func (e *Error) Wrap(err error) *Error {
	newErr := Error{
		Status:       e.Status,
		Code:         e.Code,
		Message:      e.Message,
		ErrorMessage: e.ErrorMessage,
	}
	if err != nil {
		caller := zapcore.NewEntryCaller(runtime.Caller(1)).TrimmedPath()
		newErr.ErrorMessage = fmt.Sprintf("%s: %s", caller, err.Error())
	}
	return &newErr
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
