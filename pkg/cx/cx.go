package cx

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type (
	txCtx     struct{}
	noTxCtx   struct{}
	txLockCtx struct{}
	traceCtx  struct{}
)

// NewTx wrap tx in context
func NewTx(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtx{}, db)
}

func FromTx(ctx context.Context) (any, bool) {
	v := ctx.Value(txCtx{})
	return v, v != nil
}

// NewNoTx wrap no tx in context
func NewNoTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, noTxCtx{}, true)
}

func FromNoTx(ctx context.Context) bool {
	v := ctx.Value(noTxCtx{})
	return v != nil && v.(bool)
}

// NewTxLock wrap lock tx in context
func NewTxLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, txLockCtx{}, true)
}

func FromTxLock(ctx context.Context) bool {
	v := ctx.Value(txLockCtx{})
	return v != nil && v.(bool)
}

func NewTraceCx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), traceCtx{}, c.Get(echo.HeaderXRequestID))
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(traceCtx{}).(string)
	return id
}
