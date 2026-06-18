package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string, format string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if strings.EqualFold(format, "console") {
		cfg.Encoding = "console"
	} else {
		cfg.Encoding = "json"
	}

	parsedLevel := zapcore.InfoLevel
	if err := parsedLevel.UnmarshalText([]byte(strings.ToLower(level))); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(parsedLevel)
	}

	return cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
