package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/go-playground/validator"

	"git.happyxhw.cn/happyxhw/iself/api/user"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/router/components"
)

// Serve start web serve
func Serve() {
	e := newRouter()

	s := &http.Server{
		Addr:           viper.GetString("server.addr"),
		Handler:        e,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	go func() {
		log.Info("start serve", zap.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Error("shutting down the server", zap.Error(err))
	}
}

func newRouter() *echo.Echo {
	e := echo.New()
	e.Debug = viper.GetBool("server.debug")
	e.HTTPErrorHandler = em.ErrHandler(e)
	e.Validator = &em.CustomValidator{Validator: validator.New()}

	components.InitComponent()
	initGlobalMiddleware(e)

	initRouter(e)

	return e
}

func initRouter(e *echo.Echo) {
	user.InitRouter(e)
}
