package em

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
)

// IPRateLimit by ip
func IPRateLimit(store limiter.Store, rate limiter.Rate) echo.MiddlewareFunc {
	ipRateLimiter := limiter.New(store, rate)

	// 2. Return middleware handler
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// TODO: not login -> ip
			// login -> user id
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if err != nil {
				return ErrRedis.Wrap(err)
			}

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

			if limiterCtx.Reached {
				return ErrReachLimit
			}

			return next(c)
		}
	}
}
