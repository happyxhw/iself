package handler

import (
	"net/http"

	"github.com/happyxhw/iself/pkg/ex"
)

var (
	// 400
	ErrUserExists     = ex.NewError(http.StatusConflict, 140001, "email exists")                         // 需要注册的用户已经存在了
	ErrUserSignIn     = ex.NewError(http.StatusUnauthorized, 140002, "username or password incorrect")   // 用户名或密码错误
	ErrActivation     = ex.NewError(http.StatusBadRequest, 140003, "invalid or expired activation link") // 无效或已过期的激活链接
	ErrChangePassword = ex.NewError(http.StatusBadRequest, 140004, "change password")                    // 修改密码参数错误
	ErrResetPassword  = ex.NewError(http.StatusBadRequest, 140005, "reset password")                     // 重置密码错误
	ErrOauth2Source   = ex.NewError(http.StatusBadRequest, 140006, "unknown oauth2 source")              // 未知的 oauth2 认证源
	ErrOauth2State    = ex.NewError(http.StatusBadRequest, 140007, "incorrect state")                    //  oauth2 state 校验失败

	// 500
	ErrSendEmail          = ex.NewError(http.StatusConflict, 150001, "send email")                     // ErrSendEmail 发送邮件
	ErrOauth2ExchangeCode = ex.NewError(http.StatusServiceUnavailable, 150002, "oauth2 exchange code") // 获取 access token err
	ErrGetOauth2User      = ex.NewError(http.StatusServiceUnavailable, 150003, "get oauth2 user info") // 获取 oauth2 用户信息
)
