package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFmtUnit describes the behavior of fmtUnit.
func TestFmtUnit(t *testing.T) {
	assert.Equal(t, "25 °C", fmtUnit(func() string { return "25" }, "°C"))
	assert.Equal(t, "UNKNOWN", fmtUnit(func() string { return "UNKNOWN" }, "mVA"))
	assert.Equal(t, "", fmtUnit(func() string { return "" }, "MeV"))
}
