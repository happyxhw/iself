package ex

import (
	"context"

	"github.com/happyxhw/pkg/cx"
	"github.com/labstack/echo/v4"
)

func NewTraceCtx(c echo.Context) context.Context {
	return cx.NewTraceCtx(c.Request().Context(), c.Get(echo.HeaderXRequestID))
}
