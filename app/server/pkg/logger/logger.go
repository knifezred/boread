package logger

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log      *zap.Logger
	logFile  *os.File    // 跟踪日志文件句柄，用于重初始化时关闭
	logMu    sync.Mutex
)

// Init 初始化日志
// 注意: 多次调用会先关闭之前的文件句柄
func Init(level string, logFilePath string) {
	logMu.Lock()
	defer logMu.Unlock()

	// 关闭之前的文件句柄
	if logFile != nil {
		if Log != nil {
			_ = Log.Sync()
		}
		logFile.Close()
		logFile = nil
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var cores []zapcore.Core

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), parseLevel(level))
	cores = append(cores, consoleCore)

	if logFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err == nil {
			f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				logFile = f
				fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
				fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(f), parseLevel(level))
				cores = append(cores, fileCore)
			}
		}
	}

	core := zapcore.NewTee(cores...)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
