package zlog

import (
	"my-chat/internal/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var L *zap.Logger

func Init(cfg config.LogConfig) {
	writeSyncer := getLogWriter(cfg)
	encoder := getEncoder()
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer),
		level,
	)
	L = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(L)
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func getLogWriter(cfg config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func Debug(msg string, fields ...zap.Field) {
	L.Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	L.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	L.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	L.Error(msg, fields...)
}
func With(fields ...zap.Field) *zap.Logger {
	return L.With(fields...)
}
