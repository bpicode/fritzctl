package console

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIntToCheckmark unit tests IntToCheckmark.
func TestIntToCheckmark(t *testing.T) {
	assert.NotEmpty(t, IntToCheckmark(0))
	assert.NotEmpty(t, IntToCheckmark(1))
	assert.NotEqual(t, IntToCheckmark(0), IntToCheckmark(1))
}

// TestStringToCheckmark unit tests StringToCheckmark.
func TestStringToCheckmark(t *testing.T) {
	assert.NotEmpty(t, StringToCheckmark(""))
	assert.NotEmpty(t, StringToCheckmark("0"))
	assert.NotEmpty(t, StringToCheckmark("1"))

	assert.NotEqual(t, StringToCheckmark(""), StringToCheckmark("0"))
	assert.NotEqual(t, StringToCheckmark("0"), StringToCheckmark("1"))
	assert.NotEqual(t, StringToCheckmark("1"), StringToCheckmark(""))
}
