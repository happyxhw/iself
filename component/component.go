package component

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"gorm.io/gorm"

	"github.com/spf13/viper"

	"github.com/happyxhw/iself/pkg/godb"
	"github.com/happyxhw/iself/pkg/goredis"
	"github.com/happyxhw/iself/pkg/log"
	"github.com/happyxhw/iself/pkg/mailer"
	"github.com/happyxhw/iself/pkg/oauth2x"
)

var (
	db              *gorm.DB
	rdb             *redis.Client
	ma              *mailer.Mailer
	oauth2xProvider map[string]oauth2x.Oauth2x
)

// DB return db
func DB() *gorm.DB {
	return db
}

// RDB return redis client
func RDB() *redis.Client {
	return rdb
}

// Oauth2Provider oauth2 cli map
func Oauth2Provider() map[string]oauth2x.Oauth2x {
	return oauth2xProvider
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
	oauth2xProvider = make(map[string]oauth2x.Oauth2x, len(clients))
	for _, c := range clients {
		if c.Name == oauth2x.StravaSource {
			conf := oauth2.Config{
				ClientID:     c.ClientID,
				ClientSecret: c.ClientSecret,
				Endpoint:     endpoints.Strava,
				Scopes:       c.Scopes,
			}
			oauth2xProvider[c.Name] = oauth2x.NewStrava(&conf)
		}
	}
}
