package log

import (
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type logger struct {
	zap *zap.Logger
}

var (
	wrappedLogger = &logger{}
)

// Get gets global logger
func Get() *zap.Logger {
	return wrappedLogger.zap
}

// New setups Zap to the correct log level and correct output format.
func New(logFormat, logLevel string) error {
	var zapConfig zap.Config
	switch logFormat {
	case "json":
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	default:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.DisableStacktrace = true
		zapConfig.DisableCaller = true
		zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.LevelKey = "log_level"
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02T15:04:05.000Z07:00"))
	}

	// Set the logger
	switch logLevel {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	go func(config zap.Config) {

		defaultLevel := config.Level
		var elevated bool

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGUSR1)
		for s := range c {
			if s == syscall.SIGINT {
				return
			}
			elevated = !elevated

			if elevated {
				config.Level.SetLevel(zap.DebugLevel)
				logger.Info("Log level elevated to debug")
			} else {
				logger.Info("Log level restored to original configuration", zap.String("level", logLevel))
				config.Level.SetLevel(defaultLevel.Level())
			}
		}
	}(zapConfig)
	wrappedLogger.zap = logger
	return nil
}

// NameLogger named core logger
func NameLogger(name string) {
	wrappedLogger.zap = wrappedLogger.zap.Named(name)
}

// LoggerWithMetrics add hook to zap.Logger which count log levels
func LoggerWithMetrics(statsReporter tally.Scope) {
	if statsReporter != nil {
		wrappedLogger.zap = wrappedLogger.zap.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
			if entry.Level == zap.ErrorLevel {
				statsReporter.Counter("error_count").Inc(1)
			}
			return nil
		}))
	}
}

// LoggerWithMetrics add hook to zap.Logger which count log levels
func LoggerWithErrorMetrics(scope tally.Scope, extractor TagsExtractor, metricName string) {
	if scope != nil {
		wrappedLogger.zap = wrappedLogger.zap.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return NewErrorMetricsCore(core, extractor, metricName, scope)
		}))
	}
}

// NewDevelopment setups Zap to the correct log level and correct output format for development.
func NewDevelopment(logLevel string) error {
	var zapConfig zap.Config

	zapConfig = zap.NewDevelopmentConfig()
	zapConfig.DisableStacktrace = true
	zapConfig.DisableCaller = false
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Set the logger
	switch logLevel {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.LevelKey = "log_level"
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02T15:04:05.000Z07:00"))
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}
	wrappedLogger.zap = logger
	return nil
}

// NewTest setups Zap for tests.
func NewTest() (*observer.ObservedLogs, error) {
	var zapConfig zap.Config

	zapConfig = zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zapConfig.DisableStacktrace = true
	zapConfig.DisableCaller = true
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {}
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.LevelKey = "log_level"
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02T15:04:05.000Z07:00"))
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	testCore, logs := observer.New(zap.DebugLevel)
	logger = logger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return testCore
	}))
	wrappedLogger.zap = logger

	return logs, nil
}
