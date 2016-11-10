package assert

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAssertNoErrorWithNoError is a unit test.
func TestAssertNoErrorWithNoError(t *testing.T) {
	NoError(nil)
}

// TestAssertNoErrorWithError is a unit test.
func TestAssertNoErrorWithError(t *testing.T) {
	assert.Panics(t, func() {
		NoError(errors.New("sum ting wong"))
	})
}

// TestAssertHasAtLeastWithNoError is a unit test.
func TestAssertHasAtLeastWithNoError(t *testing.T) {
	StringSliceHasAtLeast([]string{"a", "b", "c"}, 0, "Should not produce an error")
	StringSliceHasAtLeast([]string{"a", "b", "c"}, 1, "Should not produce an error")
	StringSliceHasAtLeast([]string{"a", "b", "c"}, 2, "Should not produce an error")
	StringSliceHasAtLeast([]string{"a", "b", "c"}, 3, "Should not produce an error")
}

// TestAssertHasAtLeastWithError is a unit test.
func TestAssertHasAtLeastWithError(t *testing.T) {
	assert.Panics(t, func() {
		StringSliceHasAtLeast([]string{"a", "b", "c"}, 4, "Slice too small")
	})
}

// TestIsTrueIsActuallyTrue is a unit test.
func TestIsTrueIsActuallyTrue(t *testing.T) {
	assert.NotPanics(t, func() {
		IsTrue(true)
	})
}

// TestIsTrueIsActuallyFalse is a unit test.
func TestIsTrueIsActuallyFalse(t *testing.T) {
	assert.Panics(t, func() {
		IsTrue(false)
	})
}

// TestIsEqualRegular is a unit test.
func TestIsEqualRegular(t *testing.T) {
	assert.NotPanics(t, func() {
		IsEqual(1, 1)
		IsEqual(0, 0)
		IsEqual(2, 2)
		IsEqual("abc", "abc")
	})
}

// TestIsEqualPanics is a unit test.
func TestIsEqualPanics(t *testing.T) {
	assert.Panics(t, func() {
		IsEqual(0, 1)
	})
}
