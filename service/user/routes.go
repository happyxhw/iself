package user

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/happyxhw/pkg/godb"

	"github.com/happyxhw/pkg/goredis"

	"github.com/happyxhw/pkg/mailer"

	"github.com/happyxhw/iself/pkg/ex"
	"github.com/happyxhw/iself/pkg/oauth2x"
	"github.com/happyxhw/iself/repo"
	"github.com/happyxhw/iself/service/user/controller"
	"github.com/happyxhw/iself/service/user/handler"
)

// InitRouter 初始化用户路由
func InitRouter(e *echo.Echo) {
	ag := e.Group("/api/auth")
	ug := e.Group("/api/user")
	ug.Use(ex.AuthRequired())

	cacher := repo.NewCacher(goredis.DefaultRDB())
	userRepo := repo.NewUserRepo(godb.DefaultDB())
	tokenRepo := repo.NewTokenRepo(cacher)
	aesKey := viper.GetString("secure.key")

	srv := handler.NewUserSrv(
		userRepo, tokenRepo, mailer.DefaultMailer(), cacher, []byte(aesKey),
	)
	u := controller.NewUser(srv, oauth2x.Provider())

	router(ag, ug, u)
}

func router(ag, ug *echo.Group, u *controller.User) {
	ag.POST("/sign-up", u.SignUp)                 // 注册
	ag.POST("/sign-in", u.SignIn)                 // 登录
	ag.GET("/sign-out", u.SignOut)                // 退出登录
	ag.POST("/change-password", u.ChangePassword) // 更改密码
	ag.POST("/reset-password", u.ResetPassword)   // 重置密码
	ag.POST("/active", u.Active)                  // 激活
	ag.POST("/require-email", u.SendEmail)        // 发送邮件
	ag.GET("/oauth2", u.Oauth2Callback)           // oauth2 回调接口
	ag.POST("/oauth2-state", u.Oauth2SetState)    // 设置 oauth2 state 接口, 绑定 state 到当前 session

	ug.GET("", u.Info) // 用户信息
}
