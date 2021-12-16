package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/google/uuid"

	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
)

// Serve start web serve
func Serve() {
	e := newRouter()

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func newRouter() *echo.Echo {
	e := echo.New()
	e.Debug = true
	initGlobalMiddleware(e)

	e.GET("/:id/h", func(c echo.Context) error {
		return em.NewError(http.StatusBadRequest, 400000, "debug").Wrap(errors.New("debug-err"))
	})

	return e
}

func initGlobalMiddleware(e *echo.Echo) {
	webLogger := log.NewLogger(
		&log.Config{
			Level:   viper.GetString("log.web.level"),
			Path:    viper.GetString("log.web.path"),
			Encoder: viper.GetString("log.encoder"),
		},
	)
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: uuid.NewString,
	}))
	e.HTTPErrorHandler = em.ErrHandler(e)
	e.Use(em.Logger(webLogger))
}
