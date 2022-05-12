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
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/service/user/handler"
	"git.happyxhw.cn/happyxhw/iself/service/user/types"
)

type User struct {
	srv *handler.User
}

func NewUser(srv *handler.User) *User {
	return &User{
		srv: srv,
	}
}

func (u *User) Info(c echo.Context) error {
	uc := em.GetUser(c)
	dbUser, err := u.srv.Info(em.Ctx(c), uc.ID, uc.Source)
	if err != nil {
		return err
	}

	return em.OK(c, types.NewInfo(dbUser))
}

// SignUp 用户注册
func (u *User) SignUp(c echo.Context) error {
	var req types.SignUpReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.SignUp(em.Ctx(c), req.ActiveURL, &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return em.OK(c, nil)
}

// SignIn 用户登录
func (u *User) SignIn(c echo.Context) error {
	var req types.SignInReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	dbUser, err := u.srv.SignIn(em.Ctx(c), &model.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	if dbUser.Active == 0 {
		return em.OK(c, echo.Map{
			"active": false,
			"email":  dbUser.Email,
		})
	}

	// 设置session
	u.setSession(c, dbUser, req.RememberMe)
	return em.OK(c, nil)
}

// SignOut 退出登录
func (u *User) SignOut(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     viper.GetString("session.path"),
		MaxAge:   -1,
		Domain:   viper.GetString("session.domain"),
		Secure:   viper.GetBool("session.secure"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	cookie := http.Cookie{
		Name:     "_csrf",
		Value:    random.String(32),
		Path:     "/",
		Domain:   viper.GetString("session.domain"),
		MaxAge:   -1,
		Secure:   viper.GetBool("session.secure"),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)
	_ = sess.Save(c.Request(), c.Response())
	err := em.OK(c, nil)
	return err
}

// Callback oauth2 回调接口
// TODO: redirect to sign-in failed page
// TODO: redirect to home page
func (u *User) Callback(c echo.Context) error {
	var req types.Oauth2ExchangeReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	stateSource := strings.Split(req.State, "-")
	if len(stateSource) != 2 {
		return em.ErrBadRequest
	}
	req.Source = stateSource[0]
	sess, _ := session.Get("_state", c)
	url, _ := sess.Values["url"].(string)
	state, _ := sess.Values["state"].(string)
	// delete _state session
	sess.Options = &sessions.Options{
		Path:     viper.GetString("session.path"),
		Domain:   viper.GetString("session.domain"),
		MaxAge:   -1,
		Secure:   viper.GetBool("session.secure"),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	_ = sess.Save(c.Request(), c.Response())
	if state == "" || state != req.State {
		return handler.ErrOauth2State
	}

	dbUser, err := u.srv.SignInByOauth2(em.Ctx(c), req.Source, req.Code)
	if err != nil {
		return err
	}
	u.setSession(c, &model.User{ID: dbUser.SourceID, Source: dbUser.Source}, true)
	return c.Redirect(http.StatusPermanentRedirect, url)
}

// SetState 设置 oauth2 state
func (u *User) SetState(c echo.Context) error {
	var req types.SetStateReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	u.setStateSession(c, req.State, req.URL)

	return em.OK(c, nil)
}

// Active 激活注册
func (u *User) Active(c echo.Context) error {
	var req types.ActiveReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.Active(em.Ctx(c), req.Token)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// ChangePassword 修改密码, 用户在登录状态下
func (u *User) ChangePassword(c echo.Context) error {
	var req types.ChangePasswordReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	sess, _ := session.Get("session", c)
	email, _ := sess.Values["email"].(string)
	if email == "" {
		return em.ErrForbidden
	}
	err := u.srv.ChangePassword(em.Ctx(c), email, req.Old, req.New)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// ResetPassword 用户忘记密码后重置密码，非登录状态
func (u *User) ResetPassword(c echo.Context) error {
	var req types.ResetPasswordReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.ResetPassword(em.Ctx(c), req.Password, req.Token)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// SendEmail 发送激活或重置密码邮件
func (u *User) SendEmail(c echo.Context) error {
	var req types.SendEmailReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.SendMail(em.Ctx(c), req.Email, req.Type, req.URL)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// 设置session
func (u *User) setSession(c echo.Context, user *model.User, rememberMe bool) {
	sess, _ := session.Get("session", c)
	sess.Values = map[interface{}]interface{}{
		"email":  user.Email,
		"id":     user.ID,
		"source": user.Source,
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   viper.GetInt("session.max_age"),
		Domain:   viper.GetString("session.domain"),
		Secure:   viper.GetBool("session.secure"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	cookie := http.Cookie{
		Name:     "_csrf",
		Value:    random.String(32),
		Path:     "/",
		Domain:   viper.GetString("session.domain"),
		MaxAge:   viper.GetInt("session.max_age"),
		Secure:   viper.GetBool("session.secure"),
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
		Path:     viper.GetString("session.path"),
		Domain:   viper.GetString("session.domain"),
		MaxAge:   viper.GetInt("session.max_age"),
		Secure:   viper.GetBool("session.secure"),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["state"] = state
	sess.Values["url"] = url

	_ = sess.Save(c.Request(), c.Response())
}
