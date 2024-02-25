package logger

import "go.uber.org/zap"

func Error(err error) zap.Field {
	return zap.Error(err)
}
