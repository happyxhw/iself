package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/aes"
	"git.happyxhw.cn/happyxhw/iself/pkg/ex"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/pkg/oauth2x"
	"git.happyxhw.cn/happyxhw/iself/pkg/query"
	"git.happyxhw.cn/happyxhw/iself/pkg/util"
	"git.happyxhw.cn/happyxhw/iself/service/user/types"
)

type User struct {
	aesKey []byte

	ur     UserRepo
	tr     TokenRepo
	mailer Mailer
	cacher Cacher
}

func NewUserSrv(ur UserRepo, tr TokenRepo, mailer Mailer, cacher Cacher, aesKey []byte) *User {
	return &User{
		aesKey: aesKey,

		ur:     ur,
		tr:     tr,
		mailer: mailer,
		cacher: cacher,
	}
}

func (u *User) Info(ctx context.Context, id, sourceID int64, source string) (*types.User, error) {
	var user *model.User
	var err error
	if sourceID != 0 {
		user, err = u.ur.GetBySource(ctx, source, sourceID, query.Opt{})
	} else {
		user, err = u.ur.Get(ctx, id, query.Opt{})
	}
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	if user == nil {
		return nil, ex.ErrNotFound.Msg("user not found")
	}

	return types.NewUser(user), nil
}

// SignUp 用户注册
func (u *User) SignUp(ctx context.Context, req *types.SignUpReq) error {
	// 校验用户是否存在
	data, err := u.ur.GetByEmail(ctx, req.Email, query.Fields("id"))
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	if data != nil {
		return ErrUserExists
	}
	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(passwordByte),
		Status:   int(model.WaitActiveStatus),
	}
	_, err = u.ur.Create(ctx, &user)
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	err = u.sendEmail(ctx, req.Email, ActiveEmail, req.ActiveURL)
	// 忽略邮件发送错误
	if err != nil {
		userLogger.Error("send active email", zap.Error(err), log.Ctx(ctx))
	}

	return nil
}

// SignIn 用户登录
func (u *User) SignIn(ctx context.Context, req *types.SignInReq) (*types.User, error) {
	user, err := u.ur.GetByEmail(ctx, req.Email, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	if user == nil {
		return nil, ErrUserSignIn
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err == nil {
		return types.NewUser(user), nil
	}

	return nil, ErrUserSignIn
}

// SignInByOauth2 oauth2 登录
func (u *User) SignInByOauth2(ctx context.Context, source, code string, auth oauth2x.Oauth2x) (*types.User, error) {
	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return nil, ErrOauth2ExchangeCode.Wrap(err)
	}
	oauth2User, err := auth.GetUser(ctx, token)
	if err != nil {
		return nil, ErrGetOauth2User.Wrap(err)
	}
	if oauth2User == nil || oauth2User.SourceID == 0 {
		return nil, ErrGetOauth2User
	}
	// cache token
	err = u.tr.SaveToken(ctx, token, source, oauth2User.SourceID)
	if err != nil {
		return nil, ex.ErrRedis.Wrap(err)
	}

	// 检查用户是否存在
	// 	1. 用户未绑定邮箱
	email := fmt.Sprintf("%d@%s", oauth2User.SourceID, source)
	user, err := u.ur.GetByEmail(ctx, email, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	if user != nil {
		return types.NewUser(user), nil
	}
	// 	2. 用户可能已经绑定过邮箱(用户绑定邮箱后，email 字段会替换为真实的 email)
	user, err = u.ur.GetBySource(ctx, source, oauth2User.SourceID, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	if user != nil {
		return types.NewUser(user), nil
	}
	// 用户不存在, 创建, email 为: source@source_id, oauth2 用户不需要激活
	oauth2User.Email = email
	oauth2User.Status = int(model.ActivatedStatus)
	user, err = u.ur.Create(ctx, oauth2User)
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}

	return types.NewUser(user), nil
}

// ChangePassword 更新密码
func (u *User) ChangePassword(ctx context.Context, id int64, req *types.ChangePasswordReq) error {
	user, err := u.ur.Get(ctx, id, query.Fields("id", "password"))
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	if user == nil {
		return ErrChangePassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Old))
	if err != nil {
		return ErrChangePassword
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.New), bcrypt.DefaultCost)
	params := model.UserParam{
		Password: util.String(string(hashedPassword)),
	}
	_, err = u.ur.Update(ctx, id, &params)
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	return nil
}

// ResetPassword 重置密码
func (u *User) ResetPassword(ctx context.Context, req *types.ResetPasswordReq) error {
	encrypted, err := u.decryptToken(req.Token)
	if err != nil {
		return ErrResetPassword
	}
	// email|token|timestamp
	items := strings.Split(encrypted, tokenSep)
	if len(items) != 3 {
		return ErrResetPassword
	}
	email, token := items[0], items[1]
	// 校验链接中的 token
	key := fmt.Sprintf("reset:%s", email)
	oriToken, _ := u.cacher.GetString(ctx, key)
	if oriToken == "" || oriToken != token {
		return ErrResetPassword
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	params := model.UserParam{
		Password: util.String(string(hashedPassword)),
	}
	_, err = u.ur.UpdateByEmail(ctx, email, &params)
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	// 删除原来的 token
	_, _ = u.cacher.Del(ctx, key)
	return nil
}

// Active 注册激活
func (u *User) Active(ctx context.Context, token string) error {
	decrypted, err := u.decryptToken(token)
	if err != nil {
		return ErrActivation
	}
	// email token timestamp
	items := strings.Split(decrypted, tokenSep)
	if len(items) != 3 {
		return ErrActivation
	}
	email, token := items[0], items[1]
	// 校验链接中的 token
	key := fmt.Sprintf("active:%s", email)
	oriToken, _ := u.cacher.GetString(ctx, key)
	if oriToken == "" || oriToken != token {
		return ErrActivation
	}
	// 查询用户
	params := model.UserParam{
		Status: util.Int(int(model.ActivatedStatus)),
	}

	_, err = u.ur.UpdateByEmail(ctx, email, &params)
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	// 删除原来的 token
	_, _ = u.cacher.Del(ctx, key)

	return nil
}

// SendEmail 发送邮件: 激活邮件, 重置密码邮件
func (u *User) SendEmail(ctx context.Context, email, emailType, redirectURL string) error {
	// 校验用户是否存在：
	user, err := u.ur.GetByEmail(ctx, email, query.Fields("id"))
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	if user == nil {
		return nil
	}

	return u.sendEmail(ctx, email, emailType, redirectURL)
}

func (u *User) sendEmail(ctx context.Context, email, emailType, redirectURL string) error {
	freqKey, activeKey := "active_freq:%s", "active:%s"
	if emailType == ResetEmail {
		freqKey, activeKey = "reset_freq:%s", "reset:%s"
	}
	// 频率限制
	key := fmt.Sprintf(freqKey, email)
	ret, err := u.cacher.SetNX(ctx, key, nil, emailExpire)
	if err != nil {
		return ex.ErrRedis.Wrap(err)
	}
	if !ret {
		return nil
	}

	// 保存 token
	token := util.NanoID(16)
	encryptedToken, err := u.encryptToken(email, token)
	if err != nil {
		return ex.ErrInternal.Wrap(err)
	}
	key = fmt.Sprintf(activeKey, email)
	err = u.cacher.Set(ctx, key, token, tokenExpire)
	if err != nil {
		return ex.ErrRedis.Wrap(err)
	}
	subj, body := getEmailContent(redirectURL, encryptedToken, emailType)
	err = u.mailer.Send(email, subj, body)
	if err != nil {
		return ErrSendEmail
	}
	return nil
}

func (u *User) encryptToken(email, token string) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString(email)
	buf.WriteString(tokenSep)
	buf.WriteString(token)
	buf.WriteString(tokenSep)
	buf.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	encrypted, err := aes.Encrypt(buf.Bytes(), u.aesKey)
	if err != nil {
		return "", err
	}
	tokenB64 := base64.StdEncoding.EncodeToString(encrypted)
	return url.QueryEscape(tokenB64), nil
}

func (u *User) decryptToken(token string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	decrypted, err := aes.Decrypt(encrypted, u.aesKey)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func getEmailContent(redirectURL, token, emailType string) (subj, body string) {
	link := fmt.Sprintf("%s?token=%s", redirectURL, token)
	if emailType == ActiveEmail {
		subj = "To activate your account"
		body = fmt.Sprintf("Click this link to activate your account: <p><a href=%s>%s</a></p>", link, link)
	} else {
		subj = "To reset your password"
		body = fmt.Sprintf("Click this link to reset your password: <p><a href=%s>%s</a></p>", link, link)
	}

	return subj, body
}
