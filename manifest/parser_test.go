package manifest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseAllOff test the correct parsing of an example plan file.
func TestParseAllOff(t *testing.T) {
	plan, err := ParseFile("../testdata/all_off.yml")

	assert.NoError(t, err)
	assert.NotNil(t, plan)

	assert.Len(t, plan.Switches, 3)
	assert.Equal(t, false, plan.Switches[0].State)
	assert.Equal(t, false, plan.Switches[1].State)
	assert.Equal(t, false, plan.Switches[2].State)

	assert.Len(t, plan.Thermostats, 1)
	assert.Equal(t, "ThermoOne", plan.Thermostats[0].Name)
	assert.Equal(t, float64(15), plan.Thermostats[0].Temperature)
}

// TestParseAllOn test the correct parsing of an example plan file.
func TestParseAllOn(t *testing.T) {
	plan, err := ParseFile("../testdata/all_on.yml")

	assert.NoError(t, err)
	assert.NotNil(t, plan)

	assert.Len(t, plan.Switches, 3)
	assert.Len(t, plan.Thermostats, 1)

	assert.Equal(t, true, plan.Switches[0].State)
	assert.Equal(t, true, plan.Switches[1].State)
	assert.Equal(t, true, plan.Switches[2].State)
}

// TestParseNoFileFound test the error handling when a file does not exist.
func TestParseNoFileFound(t *testing.T) {
	_, err := ParseFile("/akfnsjnqgjbqg/klksnglneglkenw/ksdgnkengkl/sdgnslgnsdl")
	assert.Error(t, err)
}

type errReader struct {
}

// Read always fails.
func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("I always fail")
}

// TestParseNotReadable test the error handling when a Reader fails.
func TestParseNotReadable(t *testing.T) {
	_, err := Parse(&errReader{})
	assert.Error(t, err)
}
