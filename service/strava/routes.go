package strava

import (
	"github.com/labstack/echo/v4"

	"git.happyxhw.cn/happyxhw/iself/component"
	"git.happyxhw.cn/happyxhw/iself/pkg/ex"
	"git.happyxhw.cn/happyxhw/iself/pkg/oauth2x"
	"git.happyxhw.cn/happyxhw/iself/pkg/trans"
	"git.happyxhw.cn/happyxhw/iself/repo"
	"git.happyxhw.cn/happyxhw/iself/service/strava/controller"
	"git.happyxhw.cn/happyxhw/iself/service/strava/handler"
)

// InitRouter 初始化用户路由
func InitRouter(e *echo.Echo) {
	g := e.Group("/api/strava")

	transRepo := trans.NewTrans(component.DB())
	sr := repo.NewStravaRepo(component.DB())
	cacher := repo.NewCacher(component.RDB())
	tr := repo.NewTokenRepo(cacher)
	auth := component.Oauth2Provider()[oauth2x.StravaSource]
	srv := handler.NewStrava(sr, tr, transRepo, auth)
	s := controller.NewStrava(srv)

	router(g, s)
}

func router(g *echo.Group, s *controller.Strava) {
	g.POST("/push", s.Push)      // 注册
	g.GET("/push", s.VerifyPush) // verify push

	g.Use(ex.AuthRequired())
	g.GET("/activities/:id", s.GetActivity)
	g.GET("/activities", s.ListActivity)

	g.GET("/activities/progress", s.GetProgressStats)
	g.GET("/activities/agg", s.GetAggStats)

	g.POST("/goals", s.CreateGoal)
	g.GET("/goals", s.QueryGoal)
	g.PUT("/goals/:id", s.UpdateGoal)
	g.DELETE("/goals/:id", s.DeleteGoal)
}
