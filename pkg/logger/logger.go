package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(debug bool) error {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	config.OutputPaths = []string{"stdout", "app.log"}

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}

func Get() *zap.Logger {
	return Log
}
