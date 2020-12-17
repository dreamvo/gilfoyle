package logging

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestLogging(t *testing.T) {
	t.Run("should create a debug logger without save on disk", func(t *testing.T) {
		l, err := NewLogger(true, false)
		assert.NoError(t, err)
		defer func() { _ = l.Sync() }()

		l.Info("test", zap.Int("test_int", 1))
	})

	t.Run("should create a production logger without save on disk", func(t *testing.T) {
		l, err := NewLogger(false, false)
		assert.NoError(t, err)
		defer func() { _ = l.Sync() }()

		l.Info("test", zap.Int("test_int", 1))
	})

	t.Run("should create a debug logger with save on disk", func(t *testing.T) {
		l, err := NewLogger(true, true)
		assert.NoError(t, err)
		defer func() { _ = l.Sync() }()
		defer func() { _ = os.RemoveAll("./logs") }()

		l.Info("test", zap.Int("test_int", 1))

		folder, err := os.Stat("./logs")
		assert.NoError(t, err)

		assert.Equal(t, "logs", folder.Name())
		assert.True(t, folder.IsDir())

		logFile, err := os.Stat("./logs/" + LogFileName)
		assert.NoError(t, err)

		assert.Equal(t, "server.log", logFile.Name())
		assert.False(t, logFile.IsDir())
	})
}
