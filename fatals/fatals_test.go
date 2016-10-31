package fatals

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAssertNoErrorWithNoError is a unit test.
func TestAssertNoErrorWithNoError(t *testing.T) {
	AssertNoError(nil)
}

// TestAssertNoErrorWithError is a unit test.
func TestAssertNoErrorWithError(t *testing.T) {
	assert.Panics(t, func() {
		AssertNoError(errors.New("sum ting wong"))
	})
}

// TestAssertHasAtLeastWithNoError is a unit test.
func TestAssertHasAtLeastWithNoError(t *testing.T) {
	AssertStringSliceHasAtLeast([]string{"a", "b", "c"}, 0, "Should not produce an error")
	AssertStringSliceHasAtLeast([]string{"a", "b", "c"}, 1, "Should not produce an error")
	AssertStringSliceHasAtLeast([]string{"a", "b", "c"}, 2, "Should not produce an error")
	AssertStringSliceHasAtLeast([]string{"a", "b", "c"}, 3, "Should not produce an error")
}

// TestAssertHasAtLeastWithError is a unit test.
func TestAssertHasAtLeastWithError(t *testing.T) {
	assert.Panics(t, func() {
		AssertStringSliceHasAtLeast([]string{"a", "b", "c"}, 4, "Slice too small")
	})
}
