package util

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugaredLogger *zap.SugaredLogger

func determineLogLevel() zap.AtomicLevel{
	flags := GetFlags()
	verbose, err := flags.GetBool("verbose")
	if nil != err {
		log.Fatal(err)
	}
	lvl := zap.NewAtomicLevel()
	if verbose {
		lvl.SetLevel(zap.DebugLevel)
	} else {
		lvl.SetLevel(zap.InfoLevel)
	}
	return lvl
}

func CreateSugaredLogger() *zap.SugaredLogger {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		determineLogLevel(),
	))
	SugaredLogger = logger.Sugar()
	return SugaredLogger
}
