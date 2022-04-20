package log

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultMaxAge     = time.Hour * 24 * 30 // 30 days
	defaultRotateTime = time.Hour * 24      // 1 day
)

// Config for log
type Config struct {
	Level      string
	Path       string
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
		EncodeTime:     zapcore.RFC3339TimeEncoder,
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
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var cores []zapcore.Core
	// 输出到文件，按天分割，warn级别及以上的会把err日志单独输出到 _err.log
	if c.Path != "" {
		if !strings.HasSuffix(c.Path, ".log") {
			c.Path += ".log"
		}
		w := writer(c)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(w), level))
		if zapLevel == zap.InfoLevel || zapLevel == zap.DebugLevel {
			c.Path = strings.TrimSuffix(c.Path, ".log")
			c.Path += "_err.log"
			errW := writer(c)
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(errW), zap.WarnLevel))
		}
	} else {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}
	// 输出到终端
	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...)
}

// log rotate
func writer(c *Config) io.Writer {
	if c.MaxAge == 0 {
		c.MaxAge = defaultMaxAge
	}
	if c.RotateTime == 0 {
		c.RotateTime = defaultRotateTime
	}
	opts := []rotateLogs.Option{
		rotateLogs.WithLinkName(c.Path),
		rotateLogs.WithMaxAge(c.MaxAge),
		rotateLogs.WithRotationTime(c.RotateTime),
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		opts = append(opts, rotateLogs.WithLocation(location))
	}
	hook, err := rotateLogs.New(
		c.Path+".%Y%m%d",
		opts...,
	)
	if err != nil {
		log.Fatalf("init rotatelogs err: %+v", err)
	}
	return hook
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
	reqID, _ := ctx.Value(echo.HeaderXRequestID).(string)
	return zap.String(echo.HeaderXRequestID, reqID)
}
