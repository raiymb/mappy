package logger

import (
	"go.uber.org/zap"
)

var l *zap.Logger

func Init(isDev bool) {
	if isDev {
		l, _ = zap.NewDevelopment()
	} else {
		l, _ = zap.NewProduction()
	}
}

func L() *zap.Logger {
	if l == nil {
		Init(true)
	}
	return l
}
