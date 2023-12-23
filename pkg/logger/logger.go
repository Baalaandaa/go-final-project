package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func GetLogger(debug bool) (*zap.Logger, error) {
	var err error
	var l *zap.Logger

	if debug {
		developmentCfg := zap.NewDevelopmentConfig()
		developmentCfg.Encoding = "console"
		developmentCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		developmentCfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		l, err = developmentCfg.Build()
	} else {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		stdout := zapcore.AddSync(os.Stdout)
		consoleEncoder := zapcore.NewConsoleEncoder(cfg)
		l = zap.New(zapcore.NewCore(consoleEncoder, stdout, zap.InfoLevel))
	}

	defer func() {
		_ = l.Sync()
	}()

	return l, err
}
