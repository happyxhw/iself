package service

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/happyxhw/pkg/godb"
	"github.com/happyxhw/pkg/goredis"
	"github.com/happyxhw/pkg/log"
	"github.com/happyxhw/pkg/mailer"

	"github.com/happyxhw/iself/pkg/oauth2x"
)

// Init 初始化 db/redis 等组件
func Init() {
	initDB()
	initRDB()
	initMailer()
	initOauth2Conf()
}

func initDB() {
	var dbC godb.Config
	_ = viper.UnmarshalKey("db", &dbC)
	dbC.Logger = log.NewLogger(
		&log.Config{
			Level:   viper.GetString("log.gorm.level"),
			Encoder: viper.GetString("log.encoder"),
		},
		zap.AddCallerSkip(3), zap.AddCaller())
	godb.InitDefaultDB(&dbC, godb.PgDB)
}

func initRDB() {
	var redC goredis.Config
	_ = viper.UnmarshalKey("redis", &redC)
	goredis.InitDefaultRDB(&redC)
}

func initMailer() {
	var c mailer.Config
	_ = viper.UnmarshalKey("mailer", &c)
	mailer.InitDefaultMailer(&c)
}

func initOauth2Conf() {
	var cfg []*oauth2x.ClientConfig
	_ = viper.UnmarshalKey("oauth2.client", &cfg)
	oauth2x.InitProvider(cfg)
}
