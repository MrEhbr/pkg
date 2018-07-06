package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	require.NotNil(t, wrappedLogger, "wrappedLogger should not be nil")
}

func Test_FromContext_NoFields(t *testing.T) {
	ctx := context.Background()
	logs, _ := NewTest()
	logger := FromContext(ctx)
	assert.NotNil(t, logger)
	assert.Empty(t, logs.TakeAll())
}

func Test_FromContext_Fields(t *testing.T) {
	ctx := context.Background()
	logs, _ := NewTest()

	field := zap.String("someKey", "someValue")
	ctx = ContextWithFields(ctx, field)

	logger := FromContext(ctx)

	require.NotNil(t, logger, "logger should not be nil")

	logger.Debug("test")

	require.Equal(t, 1, logs.Len(), "Fields does not contain the expected number of items")

	require.Contains(t, logs.TakeAll()[0].Context, field, "Fields does not contain expected field")
}
