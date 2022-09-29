package controller

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/spf13/viper"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/cx"
	"git.happyxhw.cn/happyxhw/iself/pkg/ex"
	"git.happyxhw.cn/happyxhw/iself/pkg/oauth2x"
	"git.happyxhw.cn/happyxhw/iself/pkg/util"
	"git.happyxhw.cn/happyxhw/iself/service/user/handler"
	"git.happyxhw.cn/happyxhw/iself/service/user/types"
)

type SessionConfig struct {
	Prefix string
	Key    string
	Path   string
	Domain string
	MaxAge int `mapstructure:"max_age"`
	Secure bool
}

type User struct {
	srv             *handler.User
	oauth2xProvider map[string]oauth2x.Oauth2x
	sessConfig      *SessionConfig

	homeURL   string
	activeURL string
	resetURL  string
}

func NewUser(srv *handler.User, oauth2xProvider map[string]oauth2x.Oauth2x) *User {
	var sessConfig SessionConfig
	_ = viper.UnmarshalKey("session", &sessConfig)
	return &User{
		sessConfig:      &sessConfig,
		srv:             srv,
		oauth2xProvider: oauth2xProvider,

		homeURL:   viper.GetString("web.home"),
		activeURL: viper.GetString("web.active_account"),
		resetURL:  viper.GetString("web.reset_password"),
	}
}

func (u *User) Info(c echo.Context) error {
	uc := ex.GetUser(c)
	r, err := u.srv.Info(cx.NewTraceCx(c), uc.ID, uc.SourceID, uc.Source)
	if err != nil {
		return err
	}

	return ex.OK(c, r)
}

// SignUp 用户注册
func (u *User) SignUp(c echo.Context) error {
	var req types.SignUpReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.SignUp(cx.NewTraceCx(c), &req)
	if err != nil {
		return err
	}

	return ex.OK(c, nil)
}

// SignIn 用户登录
func (u *User) SignIn(c echo.Context) error {
	var req types.SignInReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	user, err := u.srv.SignIn(cx.NewTraceCx(c), &req)
	if err != nil {
		return err
	}
	if user.Status == int(model.WaitActiveStatus) {
		return ex.OK(c, echo.Map{
			"active": false,
			"email":  user.Email,
		})
	}

	// 设置session
	u.setSession(c, user, req.RememberMe)
	return ex.OK(c, nil)
}

// SignOut 退出登录
func (u *User) SignOut(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     u.sessConfig.Path,
		MaxAge:   -1,
		Domain:   u.sessConfig.Domain,
		Secure:   u.sessConfig.Secure,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	cookie := http.Cookie{
		Name:     "_csrf",
		Value:    random.String(32),
		Path:     u.sessConfig.Path,
		Domain:   u.sessConfig.Domain,
		MaxAge:   -1,
		Secure:   u.sessConfig.Secure,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)
	_ = sess.Save(c.Request(), c.Response())
	err := ex.OK(c, nil)
	return err
}

// ChangePassword 修改密码, 用户在登录状态下
func (u *User) ChangePassword(c echo.Context) error {
	var req types.ChangePasswordReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	sess, _ := session.Get("session", c)
	if sess == nil {
		return ex.ErrForbidden
	}
	id, _ := sess.Values["id"].(float64)
	err := u.srv.ChangePassword(cx.NewTraceCx(c), int64(id), &req)
	if err != nil {
		return err
	}
	return ex.OK(c, nil)
}

// ResetPassword 用户忘记密码后重置密码，非登录状态
func (u *User) ResetPassword(c echo.Context) error {
	var req types.ResetPasswordReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.ResetPassword(cx.NewTraceCx(c), &req)
	if err != nil {
		return err
	}
	return ex.OK(c, nil)
}

// Active 用户激活
func (u *User) Active(c echo.Context) error {
	var req types.ActiveReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.Active(cx.NewTraceCx(c), req.Token)
	if err != nil {
		return err
	}

	return ex.OK(c, nil)
}

// SendEmail 发送激活或重置密码邮件
func (u *User) SendEmail(c echo.Context) error {
	var req types.SendEmailReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	url := u.activeURL
	if req.Type == handler.ResetEmail {
		url = u.resetURL
	}
	err := u.srv.SendEmail(cx.NewTraceCx(c), req.Email, req.Type, url)
	if err != nil {
		return err
	}
	return ex.OK(c, nil)
}

// Oauth2Callback oauth2 回调接口
func (u *User) Oauth2Callback(c echo.Context) error {
	var req types.Oauth2ExchangeReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	stateSource := strings.Split(req.State, "-")
	if len(stateSource) != 2 {
		return ex.ErrBadRequest
	}
	source := stateSource[0]
	cli, ok := u.oauth2xProvider[source]
	if !ok {
		return handler.ErrOauth2Source
	}
	// state 必须和当前 session 里面的 state 相匹配
	sess, _ := session.Get("_state", c)
	url, _ := sess.Values["url"].(string)
	state, _ := sess.Values["state"].(string)

	if state == "" || state != req.State {
		return handler.ErrOauth2State
	}

	// delete _state session
	sess.Options = &sessions.Options{
		Path:     u.sessConfig.Path,
		Domain:   u.sessConfig.Domain,
		MaxAge:   -1,
		Secure:   u.sessConfig.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	_ = sess.Save(c.Request(), c.Response())

	user, err := u.srv.SignInByOauth2(cx.NewTraceCx(c), source, req.Code, cli)
	if err != nil {
		return err
	}
	u.setSession(c, &types.User{SourceID: user.SourceID, Source: user.Source, Email: user.Email}, true)

	return c.Redirect(http.StatusPermanentRedirect, url)
}

// Oauth2SetState 设置 oauth2 state
func (u *User) Oauth2SetState(c echo.Context) error {
	var req types.SetOauth2StateReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	u.setStateSession(c, req.State, u.homeURL)

	return ex.OK(c, nil)
}

// 设置session
func (u *User) setSession(c echo.Context, user *types.User, rememberMe bool) {
	sess, _ := session.Get("session", c)
	sess.Values = map[interface{}]interface{}{
		"email":     user.Email,
		"id":        user.ID,
		"source_id": user.SourceID,
		"source":    user.Source,
		"status":    user.Status,
	}
	sess.Options = &sessions.Options{
		Path:     u.sessConfig.Path,
		MaxAge:   u.sessConfig.MaxAge,
		Domain:   u.sessConfig.Domain,
		Secure:   u.sessConfig.Secure,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	cookie := http.Cookie{
		Name:     "_csrf",
		Value:    util.NanoID(24),
		Path:     u.sessConfig.Path,
		Domain:   u.sessConfig.Domain,
		MaxAge:   u.sessConfig.MaxAge,
		Secure:   u.sessConfig.Secure,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	if !rememberMe {
		sess.Options.MaxAge = 0
	}
	c.SetCookie(&cookie)
	_ = sess.Save(c.Request(), c.Response())
}

// oauth2 登录, 绑定随机字符串到当前会话, 防止 csrf
func (u *User) setStateSession(c echo.Context, state, url string) {
	sess, _ := session.Get("_state", c)
	sess.Options = &sessions.Options{
		Path:     u.sessConfig.Path,
		Domain:   u.sessConfig.Domain,
		MaxAge:   0,
		Secure:   u.sessConfig.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["state"] = state
	sess.Values["url"] = url

	_ = sess.Save(c.Request(), c.Response())
}
