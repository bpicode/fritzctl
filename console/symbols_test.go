package console

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

// TestTransformToCheckmarks test the mapping of values to
// printable symbols.
func TestTransformToCheckmarks(t *testing.T) {
	testData := []struct {
		checkmark string
	}{
		{IntToCheckmark(0)},
		{IntToCheckmark(1)},
		{StringToCheckmark("")},
		{StringToCheckmark("0")},
		{StringToCheckmark("1")},
	}
	for i, testCase := range testData {
		t.Run(fmt.Sprintf("Test checkmark %d", i), func(t *testing.T) {
			assert.NotNil(t, testCase.checkmark)
			assert.NotEmpty(t, testCase.checkmark)
		})
	}
}

// TestIntToCheckmarkDisjoint test that all target symbols are different (int version).
func TestIntToCheckmarkDisjoint(t *testing.T) {
	assert.NotEqual(t, IntToCheckmark(0), IntToCheckmark(1))
}

// TestStringToCheckmarkDisjoint test that all target symbols are different (string version).
func TestStringToCheckmarkDisjoint(t *testing.T) {
	assert.NotEqual(t, StringToCheckmark(""), StringToCheckmark("0"))
	assert.NotEqual(t, StringToCheckmark("0"), StringToCheckmark("1"))
	assert.NotEqual(t, StringToCheckmark("1"), StringToCheckmark(""))
}
