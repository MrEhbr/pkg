package log

import (
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TagsExtractor func(err error) map[string]string

type core struct {
	reporter   tally.Scope
	core       zapcore.Core
	extractor  TagsExtractor
	metricName string
}

func NewErrorMetricsCore(originalCore zapcore.Core, extractor TagsExtractor, metricName string, reporter tally.Scope) zapcore.Core {
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
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
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
		c.reporter.Tagged(tags).Counter(c.metricName).Inc(1)
	}
	return c.core.Write(entry, fields)
}

func (c *core) Sync() error {
	return c.core.Sync()
}
