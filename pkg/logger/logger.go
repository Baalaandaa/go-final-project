package logger

import (
	"go.uber.org/zap"
)

func GetLogger(debug bool) (*zap.Logger, error) {
	var err error
	var l *zap.Logger
	if debug {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = l.Sync()
	}()

	return l, err
}
