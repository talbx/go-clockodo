package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugaredLogger *zap.SugaredLogger

func CreateSugaredLogger(level zapcore.Level) *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	SugaredLogger = logger.Sugar()
	return SugaredLogger
}
