package log

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/happyxhw/iself/pkg/cx"
)

const (
	defaultMaxAge     = time.Hour * 24 * 30 // 30 days
	defaultRotateTime = time.Hour * 24      // 1 day
)

// Config for log
type Config struct {
	Level      string
	Encoder    string
	MaxAge     time.Duration `mapstructure:"max_age"`
	RotateTime time.Duration `mapstructure:"rotate_time"`
}

var defaultConfig = &Config{
	Level:      "info",
	Encoder:    "console",
	MaxAge:     defaultMaxAge,
	RotateTime: defaultRotateTime,
}

// default appLogger
// no file, console type, with caller
var appLogger = NewLogger(defaultConfig, zap.AddCallerSkip(1), zap.AddCaller())

// InitAppLogger init default log
func InitAppLogger(c *Config, opts ...zap.Option) {
	appLogger = NewLogger(c, opts...)
}

// NewLogger return a new logger
func NewLogger(c *Config, opts ...zap.Option) *zap.Logger {
	if c == nil {
		c = defaultConfig
	}
	var encoder zapcore.Encoder
	// encoderConfig 编码控制
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	level := zap.NewAtomicLevel()
	zapLevel := toLevel(c.Level)
	level.SetLevel(zapLevel)
	if c.Encoder == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level), opts...)
}

// level string to zap level
func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zap.InfoLevel
}

// Debug log
func Debug(msg string, fields ...zap.Field) {
	appLogger.Debug(msg, fields...)
}

// Info log
func Info(msg string, fields ...zap.Field) {
	appLogger.Info(msg, fields...)
}

// Warn log
func Warn(msg string, fields ...zap.Field) {
	appLogger.Warn(msg, fields...)
}

// Error log
func Error(msg string, fields ...zap.Field) {
	appLogger.Error(msg, fields...)
}

// Fatal log
func Fatal(msg string, fields ...zap.Field) {
	appLogger.Fatal(msg, fields...)
}

// Panic log
func Panic(msg string, fields ...zap.Field) {
	appLogger.Panic(msg, fields...)
}

// GetLogger return appLogger
func GetLogger() *zap.Logger {
	return appLogger
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() {
	_ = appLogger.Sync()
}

func Ctx(ctx context.Context) zap.Field {
	return zap.String(echo.HeaderXRequestID, cx.RequestID(ctx))
}
