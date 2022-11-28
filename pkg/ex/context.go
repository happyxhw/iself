package ex

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/happyxhw/pkg/cx"
)

func NewTraceCtx(c echo.Context) context.Context {
	return cx.NewTraceCtx(c.Request().Context(), c.Get(echo.HeaderXRequestID))
}
