package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"

	"go.uber.org/zap"

	"github.com/happyxhw/iself/component"

	"github.com/happyxhw/iself/pkg/ex"
	"github.com/happyxhw/iself/pkg/log"
	"github.com/happyxhw/iself/third_party"
)

func initGlobalMiddleware(e *echo.Echo) {
	apiLogger := log.NewLogger(
		&log.Config{
			Level:   viper.GetString("log.web.level"),
			Encoder: viper.GetString("log.encoder"),
		},
	)
	// recovery
	if !e.Debug {
		e.Use(ex.Recover())
	}
	// request id
	e.Use(ex.RequestID())
	// access log
	e.Use(ex.Logger(apiLogger))

	initSecure(e)
	initRateLimiter(e)
	initPrometheus(e)
	initSession(e)
	initCsrf(e)
}

func initSession(e *echo.Echo) {
	rdbStore := ex.NewStore(component.RDB(), viper.GetString("session.prefix"), []byte(viper.GetString("session.key")))
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store: rdbStore,
	}))
}

func initCsrf(e *echo.Echo) {
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-TOKEN",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieDomain:   viper.GetString("session.domain"),
		CookiePath:     viper.GetString("session.path"),
		CookieSecure:   viper.GetBool("session.secure"),
		CookieHTTPOnly: false,
		CookieSameSite: http.SameSiteLaxMode,
		Skipper: func(c echo.Context) bool {
			if e.Debug {
				return true
			}
			if strings.HasPrefix(c.Path(), "/api/auth") {
				return true
			}
			return strings.HasPrefix(c.Path(), "/api/strava/push")
		},
	}))
}

func initRateLimiter(e *echo.Echo) {
	store, err := sredis.NewStoreWithOptions(component.RDB(), limiter.StoreOptions{
		Prefix: viper.GetString("ratelimit.prefix"),
	})
	if err != nil {
		log.Fatal("init limiter", zap.Error(err))
	}
	rate := limiter.Rate{
		Period: time.Duration(viper.GetInt("ratelimit.period")) * time.Second,
		Limit:  viper.GetInt64("ratelimit.limit"),
	}
	e.Use(ex.IPRateLimit(store, rate))
}

func initSecure(e *echo.Echo) {
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}))
}

func initPrometheus(e *echo.Echo) {
	p := third_party.NewPrometheus("echo", func(_ echo.Context) bool {
		return false
	})
	p.Use(e)
}
