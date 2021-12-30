package em

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/*
错误码设计
0: 成功
-1: 未知错误
4xxxx: 客户端错误
41xxx: 权限, 限制访问等
42xxx: 参数错误,
43xxx: 资源不存在等
5xxxx: 内部错误
51xxx: 基础组件错误
6xxxx: 第三方服务错误
*/

// global err
var (
	ErrBadRequest = NewError(http.StatusBadRequest, 40000, "bad request")
	ErrAuth       = NewError(http.StatusUnauthorized, 40001, "unauthorized")
	ErrForbidden  = NewError(http.StatusForbidden, 40002, "forbidden")
	ErrReachLimit = NewError(http.StatusTooManyRequests, 40003, "too many requests")
	ErrParam      = NewError(http.StatusBadRequest, 40005, "invalid parameters")
	ErrNotFound   = NewError(http.StatusNotFound, 40006, "resource not found")
	ErrConflict   = NewError(http.StatusConflict, 40007, "resource exists")

	ErrInternal = NewError(http.StatusInternalServerError, 50000, "internal server error")
	ErrDB       = NewError(http.StatusInternalServerError, 50001, "db error")
	ErrRedis    = NewError(http.StatusInternalServerError, 50002, "redis error")

	ErrThirdAPI = NewError(http.StatusServiceUnavailable, 60000, "third api error")
)

// OK 正常返回
func OK(c echo.Context, data interface{}) error {
	if data == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, data)
}

func Bind(c echo.Context, target interface{}) *Error {
	if err := c.Bind(target); err != nil {
		return ErrParam.Wrap(err)
	}
	if err := c.Validate(target); err != nil {
		return ErrParam.Wrap(err)
	}
	return nil
}
