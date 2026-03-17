package logger

import "go.uber.org/zap"

func Init(env string) *zap.Logger {
	if env == "production" {
		Log, _ := zap.NewProduction()
		return Log
	}
	Log, _ := zap.NewDevelopment()
	return Log
}
