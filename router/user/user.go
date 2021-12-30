package user

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/spf13/viper"

	"git.happyxhw.cn/happyxhw/iself/api"
	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/service"
)

type user struct {
	srv *service.User
}

func (u *user) Info(c echo.Context) error {
	userID, _ := c.Get("id").(int64)
	source, _ := c.Get("source").(string)
	dbUser, err := u.srv.Info(c.Request().Context(), userID, source)
	if err != nil {
		return err
	}

	return em.OK(c, api.NewInfo(dbUser))
}

// SignUp 用户注册
func (u *user) SignUp(c echo.Context) error {
	var req api.SignUpReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.SignUp(c.Request().Context(), req.ActiveURL, &model.User{
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
func (u *user) SignIn(c echo.Context) error {
	var req api.SignInReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	dbUser, err := u.srv.SignIn(c.Request().Context(), &model.User{
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
	u.setSession(c, dbUser)
	return em.OK(c, nil)
}

// SignOut 退出登录
func (u *user) SignOut(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     viper.GetString("session.path"),
		MaxAge:   0,
		Domain:   viper.GetString("session.domain"),
		Secure:   viper.GetBool("session.secure"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	_ = sess.Save(c.Request(), c.Response())
	err := em.OK(c, nil)
	return err
}

// Callback oauth2 回调接口
// TODO: redirect to sign-in failed page
// TODO: redirect to home page
func (u *user) Callback(c echo.Context) error {
	var req api.Oauth2ExchangeReq
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
		MaxAge:   0,
		Secure:   viper.GetBool("session.secure"),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	_ = sess.Save(c.Request(), c.Response())
	if state == "" || state != req.State {
		return service.ErrOauth2State
	}

	dbUser, err := u.srv.SignInByOauth2(c.Request().Context(), req.Source, req.Code)
	if err != nil {
		return err
	}
	u.setSession(c, &model.User{ID: dbUser.SourceID, Source: dbUser.Source})
	return c.Redirect(http.StatusPermanentRedirect, url)
}

// SetState 设置 oauth2 state
func (u *user) SetState(c echo.Context) error {
	var req api.SetStateReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	u.setStateSession(c, req.State, req.URL)

	return em.OK(c, nil)
}

// Active 激活注册
func (u *user) Active(c echo.Context) error {
	var req api.ActiveReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.Active(c.Request().Context(), req.Token)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// ChangePassword 修改密码, 用户在登录状态下
func (u *user) ChangePassword(c echo.Context) error {
	var req api.ChangePasswordReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	sess, _ := session.Get("session", c)
	email, _ := sess.Values["email"].(string)
	if email == "" {
		return em.ErrForbidden
	}
	err := u.srv.ChangePassword(c.Request().Context(), email, req.Old, req.New)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// ResetPassword 用户忘记密码后重置密码，非登录状态
func (u *user) ResetPassword(c echo.Context) error {
	var req api.ResetPasswordReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.ResetPassword(c.Request().Context(), req.Password, req.Token)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// SendEmail 发送激活或重置密码邮件
func (u *user) SendEmail(c echo.Context) error {
	var req api.SendEmailReq
	if err := em.Bind(c, &req); err != nil {
		return err
	}
	err := u.srv.SendMail(c.Request().Context(), req.Email, req.Type, req.URL)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

// 设置session
func (u *user) setSession(c echo.Context, user *model.User) {
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
	c.SetCookie(&cookie)
	_ = sess.Save(c.Request(), c.Response())
}

// oauth2 登录, 绑定随机字符串到当前会话, 防止 csrf
func (u *user) setStateSession(c echo.Context, state, url string) {
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
