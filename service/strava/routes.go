package strava

import (
	"github.com/labstack/echo/v4"

	"github.com/happyxhw/iself/component"
	"github.com/happyxhw/iself/pkg/ex"
	"github.com/happyxhw/iself/pkg/oauth2x"
	"github.com/happyxhw/iself/pkg/trans"
	"github.com/happyxhw/iself/repo"
	"github.com/happyxhw/iself/service/strava/controller"
	"github.com/happyxhw/iself/service/strava/handler"
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
