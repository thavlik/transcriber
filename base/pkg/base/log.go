package base

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var atom zap.AtomicLevel = zap.NewAtomicLevel()
var DefaultLog *zap.Logger

func Elapsed(since time.Time) zap.Field {
	return zap.String(
		"elapsed",
		time.Since(since).
			Round(time.Millisecond).
			String(),
	)
}

func init() {
	SetLogLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeDuration = zapcore.MillisDurationEncoder
	DefaultLog = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
}

func SetLogLevel() {
	if debug, ok := os.LookupEnv("DEBUG"); ok && debug != "0" {
		atom.SetLevel(zap.DebugLevel)
		return
	}
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return
	}
	switch logLevel {
	case "info":
		atom.SetLevel(zap.InfoLevel)
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "warn":
		atom.SetLevel(zap.WarnLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	case "dpanic":
		atom.SetLevel(zap.DPanicLevel)
	case "panic":
		atom.SetLevel(zap.PanicLevel)
	case "fatal":
		atom.SetLevel(zap.FatalLevel)
	default:
		panic(fmt.Sprintf("unknown log level '%s'", logLevel))
	}
}
