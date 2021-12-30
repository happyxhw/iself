package strava

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"git.happyxhw.cn/happyxhw/iself/component"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/service"
)

// InitRouter 初始化用户路由
func InitRouter(e *echo.Echo) {
	g := e.Group("/api/strava")
	stravaSrvOpt := service.StravaOption{
		DB:         component.DB(),
		RDB:        component.RDB(),
		Oauth2Conf: component.Oauth2Conf()["starava"],
		AesKey:     viper.GetString("secure.key"),
	}
	srv := service.NewStrava(&stravaSrvOpt)
	goalSrv := service.NewGoal(component.DB())
	s := strava{
		srv:     srv,
		goalSrv: goalSrv,
	}
	router(g, &s)
}

func router(g *echo.Group, s *strava) {
	g.POST("/push", s.Push)      // 注册
	g.GET("/push", s.VerifyPush) // verify push

	g.Use(em.AuthRequired())
	g.GET("/activities", s.ActivityList)
	g.GET("/activities/:id", s.Activity)
	g.GET("/activities/summary_stats", s.SummaryStatsNow)

	g.POST("/goals", s.CreateGoal)
	g.PUT("/goals", s.UpdateGoal)
	g.DELETE("/goals", s.DeleteGoal)
	g.GET("/goals", s.QueryGoal)
}
