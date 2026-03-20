package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init(env string) {
	if env == "production" {
		Log, _ = zap.NewProduction()
	} else {
		Log, _ = zap.NewDevelopment()
	}
}
