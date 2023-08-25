package util

import (
	"log"
	"log/slog"
	"os"

	"go.uber.org/zap"
)

var SugaredLogger *zap.SugaredLogger

func determineLogLevel() slog.Level {
	flags := GetFlags()
	verbose, err := flags.GetBool("verbose")
	if nil != err {
		log.Fatal(err)
	}
	if verbose {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func ConfigureLogger() {
	loggingLevel := new(slog.LevelVar)
	level := determineLogLevel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel}))
	slog.SetDefault(logger)
	loggingLevel.Set(level)
}
