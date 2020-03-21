package tests

import (
	"testing"

	"github.com/logrhythm/siem/internal/pkg/lrlog"
	"github.com/logrhythm/siem/internal/pkg/testing/capture"
	"github.com/stretchr/testify/assert"
)

func TestLogLevels(t *testing.T) {
	t.Run("test info", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.Info("This is an info log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[INFO]")
		assert.Contains(t, result, "This is an info log")
	})
	t.Run("test warning", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.Warning("This is a warning log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[WARNING]")
		assert.Contains(t, result, "This is a warning log")
	})
	t.Run("test error", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.Error("This is an error log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[ERROR]")
		assert.Contains(t, result, "This is an error log")
	})

	t.Run("test infof", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.Infof("This is an info log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[INFO]")
		assert.Contains(t, result, "This is an info log")
	})
	t.Run("test warningf", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.Warningf("This is a warning log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[WARNING]")
		assert.Contains(t, result, "This is a warning log")
	})
	t.Run("test errorf", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.Errorf("This is an error log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[ERROR]")
		assert.Contains(t, result, "This is an error log")
	})
}

func TestVLogging(t *testing.T) {
	t.Run("test v info", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.V(0).Info("This is an info log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[INFO]")
		assert.Contains(t, result, "This is an info log")
	})
	t.Run("test v warning", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.V(0).Warning("This is a warning log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[WARNING]")
		assert.Contains(t, result, "This is a warning log")
	})
	t.Run("test v error", func(t *testing.T) {
		result, err := capture.Stderr(func() {
			lrlog.SetLogWriter()
			lrlog.V(0).Error("This is an error log")
		})
		assert.Nil(t, err)
		assert.Contains(t, result, "[ERROR]")
		assert.Contains(t, result, "This is an error log")
	})
}

func TestVerbose(t *testing.T) {
	t.Run("verbosity level 0", func(t *testing.T) {
		lrlog.SetVerbosity(0)
		assert.True(t, bool(lrlog.V(0)))
		assert.False(t, bool(lrlog.V(1)))
		assert.False(t, bool(lrlog.V(2)))
		assert.False(t, bool(lrlog.V(3)))
	})
	t.Run("verbosity level 1", func(t *testing.T) {
		lrlog.SetVerbosity(1)
		assert.True(t, bool(lrlog.V(0)))
		assert.True(t, bool(lrlog.V(1)))
		assert.False(t, bool(lrlog.V(2)))
		assert.False(t, bool(lrlog.V(3)))
	})
	t.Run("verbosity level 2", func(t *testing.T) {
		lrlog.SetVerbosity(2)
		assert.True(t, bool(lrlog.V(0)))
		assert.True(t, bool(lrlog.V(1)))
		assert.True(t, bool(lrlog.V(2)))
		assert.False(t, bool(lrlog.V(3)))
	})
	t.Run("verbosity level 3", func(t *testing.T) {
		lrlog.SetVerbosity(3)
		assert.True(t, bool(lrlog.V(0)))
		assert.True(t, bool(lrlog.V(1)))
		assert.True(t, bool(lrlog.V(2)))
		assert.True(t, bool(lrlog.V(3)))
	})
}
