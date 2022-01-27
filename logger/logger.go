package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	sugaredLogger *zap.SugaredLogger
	once          sync.Once
)

func Get() *zap.SugaredLogger {
	once.Do(func() {
		sugaredLogger = NewSugaredLogger()
	})

	return sugaredLogger
}

func NewSugaredLogger() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"spotifip.log"}
	logger, _ := config.Build()
	defer func() {
		_ = logger.Sync()
	}() // flushes buffer, if any

	return logger.Sugar()
}
