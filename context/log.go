package context

import (
	"go.uber.org/zap"
)

var dlog *zap.SugaredLogger
var LogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func Debug() {
	LogLevel.SetLevel(zap.DebugLevel)
}

func init() {

	zcfg := zap.Config {
		Level:			LogLevel,
		Development:		true,
		DisableStacktrace:	true,
		Encoding:		"console",
		EncoderConfig:		zap.NewDevelopmentEncoderConfig(),
		OutputPaths:		[]string{"stderr"},
		ErrorOutputPaths:	[]string{"stderr"},
	}

	logger, _ := zcfg.Build()
	dlog = logger.Sugar()
}

func Log(desc string) *zap.SugaredLogger {
	return dlog.With(zap.String("d", desc))
}
