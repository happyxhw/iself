package godb

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
}

// NewMysqlDB return mysql db
func NewMysqlDB(dbConfig *Config) (*gorm.DB, error) {
	DB, err := createConnection(dbConfig, MysqlDB)
	return DB, err
}

// NewPgDB return postgresql db
func NewPgDB(dbConfig *Config) (*gorm.DB, error) {
	DB, err := createConnection(dbConfig, PgDB)
	return DB, err
}

// create db connection
func createConnection(dbConfig *Config, dbType dBType) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	host := dbConfig.Host
	user := dbConfig.User
	dbName := dbConfig.DB
	password := dbConfig.Password
	port := dbConfig.Port
	if host == "" {
		host = "127.0.0.1"
	}

	c := gorm.Config{
		PrepareStmt: true,
		QueryFields: true,
	}
	if dbConfig.Logger != nil {
		slowThreshold := time.Duration(dbConfig.SlowThreshold) * time.Millisecond
		c.Logger = newLogger(dbConfig.Logger, dbConfig.Level, slowThreshold, dbConfig.SQLLenThreshold)
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
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTime) * time.Second)

	return db, nil
}
