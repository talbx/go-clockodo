package util

import (
	"go.uber.org/zap"
)

var SugaredLogger *zap.SugaredLogger

func CreateSugaredLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	SugaredLogger = logger.Sugar()
	return SugaredLogger
}
