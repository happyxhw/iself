package components

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"github.com/spf13/viper"

	"git.happyxhw.cn/happyxhw/iself/pkg/godb"
	"git.happyxhw.cn/happyxhw/iself/pkg/goredis"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/pkg/mailer"
)

var (
	db         *gorm.DB
	rdb        *redis.Client
	ma         *mailer.Mailer
	oauth2Conf map[string]*oauth2.Config
)

// DB return db
func DB() *gorm.DB {
	return db
}

// RDB return redis client
func RDB() *redis.Client {
	return rdb
}

// Oauth2Conf oauth2 conf map
func Oauth2Conf() map[string]*oauth2.Config {
	return oauth2Conf
}

// Mailer send email
func Mailer() *mailer.Mailer {
	return ma
}

// InitComponent 初始化 db/redis 等组件
func InitComponent() {
	initDB()
	initRDB()
	initOauth2Conf()
	initMailer()
}

func initDB() {
	var err error
	var dbC godb.Config
	_ = viper.UnmarshalKey("db", &dbC)
	dbC.Logger = log.NewLogger(
		&log.Config{
			Level:   viper.GetString("log.gorm.level"),
			Path:    viper.GetString("log.gorm.path"),
			Encoder: viper.GetString("log.encoder"),
		},
		zap.AddCallerSkip(3), zap.AddCaller())
	db, err = godb.NewPgDB(&dbC)
	if err != nil {
		log.Fatal("init db", zap.Error(err))
	}
}

func initRDB() {
	var err error
	var redC goredis.Config
	_ = viper.UnmarshalKey("redis", &redC)
	rdb, err = goredis.NewRedis(&redC)
	if err != nil {
		log.Fatal("init redis", zap.Error(err))
	}
}

func initMailer() {
	var c mailer.Config
	_ = viper.UnmarshalKey("mailer", &c)
	ma = mailer.NewMailer(&c)
}

type oauth2Client struct {
	Name         string
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	Scopes       []string
}

func initOauth2Conf() {
	var clients []*oauth2Client
	_ = viper.UnmarshalKey("oauth2.client", &clients)
	oauth2Conf = make(map[string]*oauth2.Config, len(clients))
	for _, c := range clients {
		var endpoint oauth2.Endpoint
		switch c.Name {
		case "github":
			endpoint = github.Endpoint
		case "google":
			endpoint = google.Endpoint
		default:
			continue
		}
		oauth2Conf[c.Name] = &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			Endpoint:     endpoint,
			Scopes:       c.Scopes,
		}
	}
}
