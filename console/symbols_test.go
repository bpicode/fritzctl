package console

import (
	"fmt"
	"testing"

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

// TestStringToCheckmarkInverse test the inverse relations of checkmarks.
func TestStringToCheckmarkInverse(t *testing.T) {
	assert.Equal(t, Stoc("1"), Stoc("0").Inverse())
	assert.Equal(t, Stoc(""), Stoc("").Inverse())
	assert.Equal(t, Stoc("0"), Stoc("1").Inverse())
	assert.Equal(t, Stoc("0"), Stoc("0").Inverse().Inverse())
	assert.Equal(t, Stoc("1"), Stoc("1").Inverse().Inverse())
	assert.Equal(t, Stoc(""), Stoc("").Inverse().Inverse())
}

// TestBooleanToCheckmark test the checkmark creation of a bool.
func TestBooleanToCheckmark(t *testing.T) {
	assert.NotZero(t, Btoc(true))
	assert.NotZero(t, Btoc(false))
}
