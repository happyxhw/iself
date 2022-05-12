package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"git.happyxhw.cn/happyxhw/iself/component"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	stravaRouter "git.happyxhw.cn/happyxhw/iself/service/strava"
	userRouter "git.happyxhw.cn/happyxhw/iself/service/user"
	weatherRouter "git.happyxhw.cn/happyxhw/iself/service/weather"
)

// Serve start web serve
func Serve() {
	e := newRouter()

	s := &http.Server{
		Addr:           viper.GetString("server.addr"),
		Handler:        e,
		MaxHeaderBytes: 1 << 20,
	}

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
	e.Validator = em.NewValidator()
	e.IPExtractor = echo.ExtractIPFromRealIPHeader(echo.TrustLinkLocal(true), echo.TrustPrivateNet(true))

	component.InitComponent()
	initGlobalMiddleware(e)

	initRouter(e)

	return e
}

func initRouter(e *echo.Echo) {
	userRouter.InitRouter(e)
	stravaRouter.InitRouter(e)
	weatherRouter.InitRouter(e)
}
