package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNonPanicLoggers ensures that informational loggers do not panic.
func TestNonPanicLoggers(t *testing.T) {
	assert.NotPanics(t, func() {
		Debug("some", "debug", "message")
		Info("some", "info", "message")
		Success("some", "success", "message")
		Warn("some", "warning", "message")
		Error("some", "error", "message")
	})
}

// TestPanicLogger ensures that panic loggers do panic.
func TestPanicLogger(t *testing.T) {
	assert.Panics(t, func() {
		Panic("ooops")
	})
}

// TestRegularLogLevelConfigs ensures that all predefined levels can be configured.
func TestRegularLogLevelConfigs(t *testing.T) {
	defer ConfigureLogLevel("info")
	for _, name := range levelNames.keys() {
		err := ConfigureLogLevel(name)
		assert.NoError(t, err)
		assert.NotNil(t, ls)
		assert.NotNil(t, ls.panicLog)
		assert.NotNil(t, ls.warn)
		assert.NotNil(t, ls.success)
		assert.NotNil(t, ls.info)
		assert.NotNil(t, ls.debug)
	}
}

// TestInvalidLogLevelConfig tests that an unknown level config returns an error.
func TestInvalidLogLevelConfig(t *testing.T) {
	defer ConfigureLogLevel("info")
	err := ConfigureLogLevel("does-not-exist")
	assert.Error(t, err)
}

// TestLookupSanity asserts that the lookup table is sane.
func TestLookupSanity(t *testing.T) {
	keys := levelNames.keys()
	assert.NotEmpty(t, keys)

	vs := make(map[*printers]interface{})
	for _, v := range levelNames {
		assert.NotEmpty(t, v)
		vs[v] = nil
	}
	assert.Equal(t, len(keys), len(vs))
}

// TestLogLevelByName asserts that the LevelByName is working properly for pre-defined log levels.
func TestLogLevelByName(t *testing.T) {
	for _, name := range levelNames.keys() {
		_, err := byName(name)
		assert.NoError(t, err)
	}
}

// TestLogLevelByNameWithInvalidArg asserts that the LevelByName is working properly for invalid loglevel argument.
func TestLogLevelByNameWithInvalidArg(t *testing.T) {
	_, err := byName("does-not-exist")
	assert.Error(t, err)
}
