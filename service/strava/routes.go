package strava

import (
	"github.com/labstack/echo/v4"

	"github.com/happyxhw/pkg/godb"

	"github.com/happyxhw/pkg/goredis"

	"github.com/happyxhw/pkg/trans"

	"github.com/happyxhw/iself/pkg/ex"
	"github.com/happyxhw/iself/pkg/oauth2x"
	"github.com/happyxhw/iself/repo"
	"github.com/happyxhw/iself/service/strava/controller"
	"github.com/happyxhw/iself/service/strava/handler"
)

// InitRouter 初始化用户路由
func InitRouter(e *echo.Echo) {
	g := e.Group("/api/strava")

	transRepo := trans.NewTrans(godb.DefaultDB())
	sr := repo.NewStravaRepo(godb.DefaultDB())
	cacher := repo.NewCacher(goredis.DefaultRDB())
	tr := repo.NewTokenRepo(cacher)
	auth := oauth2x.Provider()[oauth2x.StravaSource]
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
