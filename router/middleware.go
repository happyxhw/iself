package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"

	sredis "github.com/ulule/limiter/v3/drivers/store/redis"

	"go.uber.org/zap"

	"git.happyxhw.cn/happyxhw/iself/component"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
)

func initGlobalMiddleware(e *echo.Echo) {
	apiLogger := log.NewLogger(
		&log.Config{
			Level:   viper.GetString("log.web.level"),
			Path:    viper.GetString("log.web.path"),
			Encoder: viper.GetString("log.encoder"),
		},
	)
	// recovery
	e.Use(middleware.Recover())
	// request id
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: uuid.NewString,
	}))
	// access log
	e.Use(em.Logger(apiLogger))

	initSecure(e)
	initRateLimiter(e)
	initPrometheus(e)
	initSession(e)
	initCsrf(e)
	initCors(e)
}

func initSession(e *echo.Echo) {
	rdbStore := em.NewStore(component.RDB(), viper.GetString("session.prefix"), []byte(viper.GetString("session.key")))
	e.Use(session.MiddlewareWithConfig(session.Config{
		Store: rdbStore,
	}))
}

func initCsrf(e *echo.Echo) {
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLength:    32,
		TokenLookup:    "header:X-XSRF-TOKEN",
		ContextKey:     "csrf",
		CookieName:     "_csrf",
		CookieDomain:   viper.GetString("session.domain"),
		CookiePath:     "/",
		CookieMaxAge:   viper.GetInt("session.max_age"),
		CookieSecure:   viper.GetBool("session.secure"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteLaxMode,
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Path(), "/api/auth") {
				return true
			}
			return strings.HasPrefix(c.Path(), "/api/strava/push")
		},
	}))
}

func initCors(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost", "ifit.happyxhw.com"},
		AllowMethods: []string{http.MethodGet},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		MaxAge:       86400,
	}))
}

func initRateLimiter(e *echo.Echo) {
	store, err := sredis.NewStoreWithOptions(component.RDB(), limiter.StoreOptions{
		Prefix: "limiter",
	})
	if err != nil {
		log.Fatal("init limiter", zap.Error(err))
	}
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  120,
	}
	e.Use(em.IPRateLimit(store, rate))
}

func initSecure(e *echo.Echo) {
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}))
}

func initPrometheus(e *echo.Echo) {
	p := prometheus.NewPrometheus("echo", func(_ echo.Context) bool {
		return false
	})
	p.Use(e)
}
