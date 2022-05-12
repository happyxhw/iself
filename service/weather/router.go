package weather

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"git.happyxhw.cn/happyxhw/iself/pkg/log"
)

const appID = "82a3f058f5d7010b2a63a18dff557de7"

// InitRouter 初始化用户路由
func InitRouter(e *echo.Echo) {
	g := e.Group("/api/open_weather")
	remote, err := url.Parse(fmt.Sprintf("https://api.openweathermap.org/data/2.5?appid=%s", appID))
	if err != nil {
		log.Fatal("init proxy", zap.Error(err))
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(request *http.Request) {
		targetQuery := remote.RawQuery
		request.URL.Scheme = remote.Scheme
		request.URL.Host = remote.Host
		request.Host = remote.Host
		request.URL.Path = fmt.Sprintf("%s/%s", remote.Path, strings.TrimPrefix(request.URL.Path, "/api/open_weather/"))

		if targetQuery == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = targetQuery + request.URL.RawQuery
		} else {
			request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
		}
	}
	g.GET("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
}
