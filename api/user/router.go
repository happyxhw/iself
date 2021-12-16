package user

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/router/components"
	"git.happyxhw.cn/happyxhw/iself/service"
)

// InitRouter 初始化用户路由
func InitRouter(e *echo.Echo) {
	ag := e.Group("/api/auth")
	ug := e.Group("/api/user")
	ug.Use(em.AuthRequired())
	userSrvOpt := service.UserOption{
		DB:         components.DB(),
		RDB:        components.RDB(),
		Oauth2Conf: components.Oauth2Conf(),
		Ma:         components.Mailer(),
		AesKey:     viper.GetString("secure.key"),
	}
	srv := service.NewUser(&userSrvOpt)
	u := user{
		srv: srv,
	}
	router(ag, ug, &u)
}

func router(ag, ug *echo.Group, u *user) {
	ag.POST("/sign-up", u.SignUp)                 // 注册
	ag.POST("/sign-in", u.SignIn)                 // 登录
	ag.POST("/active", u.Active)                  // 激活
	ag.POST("/change-password", u.ChangePassword) // 重设密码
	ag.POST("/reset-password", u.ResetPassword)   // 重设密码
	ag.POST("/require-email", u.SendEmail)        // 发送邮件
	ag.GET("/oauth2", u.Callback)                 // oauth2 回调接口
	ag.POST("/oauth2_state", u.SetState)          // 设置 oauth2 state 接口, 绑定 state 到当前 session

	ug.GET("", u.Info)
}
