package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(level string, logFile string) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var cores []zapcore.Core

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), parseLevel(level))
	cores = append(cores, consoleCore)

	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0755); err == nil {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
				fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), parseLevel(level))
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