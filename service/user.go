package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/google/go-github/v41/github"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/aes"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/pkg/mailer"
	"git.happyxhw.cn/happyxhw/iself/pkg/util"
)

/*
用户数据服务
*/

var userLogger = log.GetLogger().Named("user")

var (
	// ErrUserExists 需要注册的用户已经存在了
	ErrUserExists = em.NewError(http.StatusConflict, 40101, "email exists")
	// ErrUserSignIn 用户名或密码错误
	ErrUserSignIn = em.NewError(http.StatusUnauthorized, 40102, "username or password incorrect")
	// ErrOauth2Source 未知的 oauth2 认证源
	ErrOauth2Source = em.NewError(http.StatusBadRequest, 40103, "unknown oauth2 source")
	// ErrOauth2State oauth2 state 校验失败
	ErrOauth2State = em.NewError(http.StatusBadRequest, 40104, "incorrect state")
	// ErrActive 无效或已过期的激活链接
	ErrActive = em.NewError(http.StatusBadRequest, 40105, "invalid or expired active link")
	// ErrOldPassword 老密码错误
	ErrOldPassword = em.NewError(http.StatusBadRequest, 40106, "invalid password")
	// ErrResetPassword 重置密码错误
	ErrResetPassword = em.NewError(http.StatusBadRequest, 40107, "reset password")
	// ErrOauth2ExchangeCode 获取 access token err
	ErrOauth2ExchangeCode = em.NewError(http.StatusServiceUnavailable, 60101, "oauth2 exchange code")
	// ErrGetOauth2User 获取 oauth2 用户信息
	ErrGetOauth2User = em.NewError(http.StatusServiceUnavailable, 60102, "get oauth2 user info")
)

// User 用户数据服务
type User struct {
	db         *gorm.DB
	rdb        *redis.Client
	oauth2Conf map[string]*oauth2.Config
	tokenSrv   *Token
	ma         *mailer.Mailer
	aesKey     []byte
}

// UserOption config for user
type UserOption struct {
	DB         *gorm.DB
	RDB        *redis.Client
	Oauth2Conf map[string]*oauth2.Config
	Ma         *mailer.Mailer
	AesKey     string
}

// NewUser 返回数据服务实例
func NewUser(opt *UserOption) *User {
	return &User{
		db:         opt.DB,
		rdb:        opt.RDB,
		oauth2Conf: opt.Oauth2Conf,
		tokenSrv:   &Token{rdb: opt.RDB},
		ma:         opt.Ma,
		aesKey:     []byte(opt.AesKey),
	}
}

func (u *User) Info(_ context.Context, id int64, source string) (*model.User, *em.Error) {
	var dbUser model.User
	var query *gorm.DB
	if source == "" {
		query = u.db.Select("id, email, name, avatar_url").Where("id = ?", id)
	} else {
		// TODO
		query = u.db.Select("id, avatar_url").Where("source_id = ? AND source = ?", id, source)
	}
	err := query.Find(&dbUser).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	return &dbUser, nil
}

// SignUp 注册用户
// oauth2 用户后面可通过应用内重置密码的方式设置密码
func (u *User) SignUp(ctx context.Context, redirectURL string, user *model.User) *em.Error {
	// 校验用户是否存在
	var dbUser model.User
	err := u.db.Select("id").Where("email = ?", user.Email).Find(&dbUser).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if dbUser.ID > 0 {
		return ErrUserExists
	}

	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordByte)
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if errDB := u.db.Create(&user).Error; errDB != nil {
		return em.ErrDB.Wrap(errDB)
	}
	emErr := u.SendMail(ctx, user.Email, "active", redirectURL)
	if emErr != nil {
		userLogger.Error("send active email", zap.Error(err))
	}
	return nil
}

// SignIn 用户登录
func (u *User) SignIn(_ context.Context, user *model.User) (*model.User, *em.Error) {
	var dbUser model.User
	err := u.db.Select("id, email, password, active").Where("email = ?", user.Email).Find(&dbUser).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err == nil {
		return &dbUser, nil
	}

	return nil, ErrUserSignIn
}

// SignInByOauth2 oauth2 登录
func (u *User) SignInByOauth2(ctx context.Context, source, code string) (*model.User, *em.Error) {
	conf, ok := u.oauth2Conf[source]
	if !ok {
		return nil, ErrOauth2Source
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, ErrOauth2ExchangeCode.Wrap(err)
	}
	cli := conf.Client(ctx, token)
	// TODO: 支持其他客户端
	client := github.NewClient(cli)
	authenticatedUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, ErrGetOauth2User.Wrap(err)
	}
	if authenticatedUser.ID == nil {
		return nil, ErrGetOauth2User
	}
	if saveErr := u.tokenSrv.SaveToken(token, source, *authenticatedUser.ID); saveErr != nil {
		return nil, em.ErrRedis.Wrap(err)
	}

	var uo model.UserOauth2
	err = u.db.Select("id").
		Where("source_id = ? AND source = ?", authenticatedUser.ID, source).Find(&uo).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	if uo.ID > 0 {
		return &model.User{ID: *authenticatedUser.ID, Source: source}, nil
	}

	if err := u.db.Create(newGithubUser(authenticatedUser)).Error; err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	return &model.User{ID: *authenticatedUser.ID, Source: source}, nil
}

// Active 注册激活
func (u *User) Active(ctx context.Context, token string) *em.Error {
	decrypted, err := u.decryptToken(token)
	if err != nil {
		return ErrActive
	}
	// email|token|timestamp
	items := strings.Split(decrypted, "|")
	if len(items) != 3 {
		return ErrActive
	}
	// 校验链接中的 token
	key := fmt.Sprintf("active:%s", items[0])
	oriToken, _ := u.rdb.Get(ctx, key).Result()
	if oriToken == "" || oriToken != items[1] {
		return ErrActive
	}
	updatedMap := map[string]interface{}{
		"active":     1,
		"updated_at": time.Now(),
	}
	err = u.db.Model(&model.User{}).Where("email = ?", items[0]).Updates(updatedMap).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	// 删除原来的 token
	_ = u.rdb.Del(ctx, key)
	return nil
}

// SendMail 发送注册激活、重置密码邮件
func (u *User) SendMail(ctx context.Context, email, emailType, redirectURL string) *em.Error {
	freqKey, tokenKey := "active_freq:%s", "active:%s" //nolint:gosec
	if emailType == "reset" {
		freqKey, tokenKey = "reset_freq:%s", "reset:%s" //nolint:gosec
	}
	// 限制发送频率, 2 分钟发送一次
	key := fmt.Sprintf(freqKey, email)
	ret, _ := u.rdb.SetNX(ctx, key, nil, time.Minute*2).Result()
	if !ret {
		return nil
	}
	var dbUser model.User
	var err error
	if emailType == "active" {
		// 校验用户是否存在或已经激活
		err = u.db.Select("id").Where("email = ? AND active = ?", email, 0).Find(&dbUser).Error
	} else {
		// 未激活用户也可以重置
		err = u.db.Select("id").Where("email = ? ", email).Find(&dbUser).Error
	}
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if dbUser.ID == 0 {
		return nil
	}

	// 有效期 30 分钟
	token := util.GenerateToken(16)
	key = fmt.Sprintf(tokenKey, email)
	// 失效原来的 token
	_ = u.rdb.Del(ctx, key)
	err = u.rdb.Set(ctx, key, token, time.Minute*30).Err()
	if err != nil {
		return em.ErrRedis.Wrap(err)
	}
	token, err = u.encryptToken(email, token)
	if err != nil {
		return em.ErrInternal.Wrap(err)
	}
	subj := "Welcome"
	link := fmt.Sprintf("%s?token=%s", redirectURL, token)
	body := fmt.Sprintf("Click this link to %s your account: <p><a href=%s>%s</a></p>", emailType, link, link)
	err = u.ma.DialAndSend(email, subj, body)
	if err != nil {
		return em.ErrThirdAPI.Wrap(err)
	}
	return nil
}

// ChangePassword 更新密码
func (u *User) ChangePassword(_ context.Context, email, oldPass, newPass string) *em.Error {
	var dbUser model.User
	err := u.db.Select("id, password").Where("email = ? AND active = ?", email, 1).Find(&dbUser).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if dbUser.ID > 0 {
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(oldPass))
		if err != nil {
			return ErrOldPassword
		}
		passwordByte, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		updatedMap := map[string]interface{}{
			"password":   string(passwordByte),
			"updated_at": time.Now(),
		}
		err := u.db.Model(&model.User{}).Where("email = ?", email).Updates(updatedMap).Error
		if err != nil {
			return em.ErrDB.Wrap(err)
		}
	}
	return nil
}

// ResetPassword 重置密码
func (u *User) ResetPassword(ctx context.Context, password, token string) *em.Error {
	encrypted, err := u.decryptToken(token)
	if err != nil {
		return ErrResetPassword
	}
	// email|token|timestamp
	items := strings.Split(encrypted, "|")
	if len(items) != 3 {
		return ErrResetPassword
	}
	// 校验链接中的 token
	key := fmt.Sprintf("reset:%s", items[0])
	oriToken, _ := u.rdb.Get(ctx, key).Result()
	if oriToken == "" || oriToken != items[1] {
		return ErrResetPassword
	}
	passwordByte, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	updatedMap := map[string]interface{}{
		"password":   string(passwordByte),
		"updated_at": time.Now(),
	}
	err = u.db.Model(&model.User{}).Where("email = ?", items[0]).Updates(updatedMap).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	// 删除原来的 token
	_ = u.rdb.Del(ctx, key)
	return nil
}

func (u *User) encryptToken(email, token string) (string, error) {
	buf := bytes.Buffer{}
	buf.WriteString(email)
	buf.WriteString("|")
	buf.WriteString(token)
	buf.WriteString("|")
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

func newGithubUser(u *github.User) *model.UserOauth2 {
	uo := &model.UserOauth2{
		Source:   "github",
		SourceID: *u.ID,
	}
	if u.AvatarURL != nil {
		uo.AvatarURL = *u.AvatarURL
	}

	return uo
}
