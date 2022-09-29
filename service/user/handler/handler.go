package handler

import (
	"context"
	"time"

	"golang.org/x/oauth2"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/pkg/query"
	"git.happyxhw.cn/happyxhw/iself/repo"
)

var userLogger = log.GetLogger().Named("user")

//go:generate mockgen -destination=./mocks/mock_user_repo.go -package=mocks . UserRepo
type UserRepo interface {
	Create(ctx context.Context, u *model.User) (*model.User, error)
	Get(ctx context.Context, id int64, opt query.Opt) (*model.User, error)
	GetByEmail(ctx context.Context, email string, opt query.Opt) (*model.User, error)
	GetBySource(ctx context.Context, source string, sourceID int64, opt query.Opt) (*model.User, error)
	Update(ctx context.Context, id int64, params *model.UserParam) (int64, error)
	UpdateByEmail(ctx context.Context, email string, params *model.UserParam) (int64, error)
}

//go:generate mockgen -destination=./mocks/mock_token_repo.go -package=mocks . TokenRepo
type TokenRepo interface {
	SaveToken(ctx context.Context, token *oauth2.Token, source string, sourceID int64) error
	GetToken(ctx context.Context, source string, sourceID int64, fn repo.Refresher) (*oauth2.Token, error)
}

//go:generate mockgen -destination=./mocks/mock_cacher.go -package=mocks . Cacher
type Cacher interface {
	GetString(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val any, ex time.Duration) error
	SetNX(ctx context.Context, key string, val any, ex time.Duration) (bool, error)
	Del(ctx context.Context, key string) (int64, error)
}

//go:generate mockgen -destination=./mocks/mock_oauth2x.go -package=mocks git.happyxhw.cn/happyxhw/iself/pkg/oauth2x Oauth2x

//go:generate mockgen -destination=./mocks/mock_mailer.go -package=mocks . Mailer
type Mailer interface {
	Send(to, subj, body string) error
}

const (
	emailExpire   = time.Minute * 2 // 邮件发送频率限制
	tokenExpire   = time.Minute * 30
	oauth2Timeout = time.Minute * 1
)

const (
	activeEmail = "active" // 用户激活邮件
	resetEmail  = "reset"  // 用户重置密码邮件
)

const tokenSep = " "
