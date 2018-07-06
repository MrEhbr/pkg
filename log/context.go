package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// key type is unexported to prevent collisions with context keys defined in other packages.
type key int

const (
	// loggerKey is the context key for logger
	loggerKey key = iota
	// fieldsKey is the context key for the Fields.
	fieldsKey
)

// NewContext creates a new context the given contextual fields
func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, loggerKey, FromContext(ctx).With(fields...))
}

// FromContext returns a logger from the given context
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return wrappedLogger.zap
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger.With(FieldsFromContext(ctx)...)
	}
	return wrappedLogger.zap.With(FieldsFromContext(ctx)...)
}

// FieldsFromContext retrieves the Fields from ctx.
func FieldsFromContext(ctx context.Context) []zap.Field {
	if fields, ok := ctx.Value(fieldsKey).([]zap.Field); ok {
		return fields
	}
	return []zap.Field{}
}

// ContextWithFields set fields in a new context based on ctx, and returns this
// context. Any Fields defined in ctx will be overriden.
func ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, fieldsKey, fields)
}

// AddFieldToCtx add set of fields in a new context based on ctx, and returns this context.
func AddFieldToCtx(ctx context.Context, fields ...zap.Field) context.Context {
	ctxFields := FieldsFromContext(ctx)
	ctxFields = append(ctxFields, fields...)
	return ContextWithFields(ctx, ctxFields...)
}
