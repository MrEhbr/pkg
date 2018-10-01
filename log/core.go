package log

import (
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TagsExtractor func(err error) map[string]string

type core struct {
	reporter   tally.StatsReporter
	core       zapcore.Core
	extractor  TagsExtractor
	metricName string
}

func NewErrorMetricsCore(originalCore zapcore.Core, extractor TagsExtractor, metricName string, reporter tally.StatsReporter) zapcore.Core {
	return &core{
		reporter:   reporter,
		core:       originalCore,
		extractor:  extractor,
		metricName: metricName,
	}
}

func (c *core) Enabled(level zapcore.Level) bool {
	return c.core.Enabled(level)
}

func (c *core) With(fields []zapcore.Field) zapcore.Core {
	return c.core.With(fields)
}

func (c *core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return c.core.Check(ent, ce)
}

func (c *core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if entry.Level == zap.ErrorLevel {
		var tags map[string]string
		for _, field := range fields {
			if field.Type == zapcore.ErrorType {
				tags = c.extractor(field.Interface.(error))
				break
			}
		}
		c.reporter.ReportCounter(c.metricName, tags, 1)
	}
	return c.core.Write(entry, fields)
}

func (c *core) Sync() error {
	return c.core.Sync()
}