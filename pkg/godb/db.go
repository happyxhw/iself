package godb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/happyxhw/iself/pkg/cx"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

var (
	// ErrUnknownDBType unknown db type
	ErrUnknownDBType = errors.New("unknown db type")
)

// database type
type dBType int8

const (
	// MysqlDB mysql
	MysqlDB dBType = iota
	// PgDB postgresql
	PgDB
)

// Config db config
type Config struct {
	User            string
	Password        string
	Host            string
	Port            int
	DB              string
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	MaxLifeTime     int `mapstructure:"max_life_time"`
	Logger          *zap.Logger
	Level           string
	SlowThreshold   int `mapstructure:"slow_threshold"`
	SQLLenThreshold int `mapstructure:"sql_len_threshold"`

	MetricsPort uint32 `mapstructure:"metrics_port"`
	Prometheus  bool
}

// NewMysqlDB return mysql db
func NewMysqlDB(cfg *Config) (*gorm.DB, error) {
	DB, err := createConnection(cfg, MysqlDB)
	return DB, err
}

// NewPgDB return postgresql db
func NewPgDB(cfg *Config) (*gorm.DB, error) {
	DB, err := createConnection(cfg, PgDB)
	return DB, err
}

// create db connection
func createConnection(cfg *Config, dbType dBType) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	host := cfg.Host
	user := cfg.User
	dbName := cfg.DB
	password := cfg.Password
	port := cfg.Port
	if host == "" {
		host = "127.0.0.1"
	}

	c := gorm.Config{
		PrepareStmt: true,
		QueryFields: true,
	}
	if cfg.Logger != nil {
		slowThreshold := time.Duration(cfg.SlowThreshold) * time.Millisecond
		c.Logger = newLogger(cfg.Logger, cfg.Level, slowThreshold, cfg.SQLLenThreshold)
	}
	switch dbType {
	case MysqlDB:
		if port == 0 {
			port = 3306
		}
		url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&interpolateParams=true",
			user, password, host, port, dbName)
		db, err = gorm.Open(mysql.Open(url), &c)
	case PgDB:
		if port == 0 {
			port = 5432
		}
		url := fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbName, password,
		)
		db, err = gorm.Open(postgres.Open(url), &c)
	default:
		return nil, ErrUnknownDBType
	}
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	if cfg.Prometheus {
		_ = db.WithContext(cx.NewMetricsCx(context.Background())).Use(prometheus.New(prometheus.Config{
			DBName:          cfg.DB,
			RefreshInterval: 30,
			StartServer:     true,
			HTTPServerPort:  cfg.MetricsPort,
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.Postgres{
					Prefix: "ifly_gorm_pg_",
				},
			}, // 用户自定义指标
		}))
	}

	return db, nil
}
