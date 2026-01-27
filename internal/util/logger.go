package util

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var DBLogger *zap.Logger

// LogConfig 定义日志配置接口，避免循环依赖
type LogConfig interface {
	GetLevel() string
	GetFilename() string
	GetDBFilename() string
	GetMaxSize() int
	GetMaxBackups() int
	GetMaxAge() int
	GetCompress() bool
	GetLogConsole() bool
	GetShowDBLog() bool
	GetDBLogLevel() string
}

func InitLogger(cfg LogConfig) {
	// 确保日志目录存在
	if cfg.GetFilename() != "" {
		_ = os.MkdirAll(filepath.Dir(cfg.GetFilename()), 0755)
	}
	if cfg.GetDBFilename() != "" {
		_ = os.MkdirAll(filepath.Dir(cfg.GetDBFilename()), 0755)
	}

	Logger = createZapLogger(
		cfg.GetLevel(),
		cfg.GetFilename(),
		cfg.GetMaxSize(),
		cfg.GetMaxBackups(),
		cfg.GetMaxAge(),
		cfg.GetCompress(),
		cfg.GetLogConsole(),
	)

	DBLogger = createZapLogger(
		cfg.GetDBLogLevel(),
		cfg.GetDBFilename(),
		cfg.GetMaxSize(),
		cfg.GetMaxBackups(),
		cfg.GetMaxAge(),
		cfg.GetCompress(),
		cfg.GetShowDBLog(),
	)
}

func createZapLogger(level, filename string, maxSize, maxBackups, maxAge int, compress, logConsole bool) *zap.Logger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var cores []zapcore.Core

	// 文件输出
	if filename != "" {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		})
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			zapLevel,
		))
	}

	// 控制台输出
	if logConsole {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		))
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller())
}
