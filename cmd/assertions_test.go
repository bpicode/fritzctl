package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAssertNoErrorWithNoError is a unit test.
func TestAssertNoErrorWithNoError(t *testing.T) {
	assert.NotPanics(t, func() {
		assertNoErr(nil, "would-be-context")
	})
}

// TestAssertNoErrorWithError is a unit test.
func TestAssertNoErrorWithError(t *testing.T) {
	assert.Panics(t, func() {
		assertNoErr(errors.New("we too low"), "sum ting wong")
	})
}

// TestAssertHasAtLeastWithNoError is a unit test.
func TestAssertHasAtLeastWithNoError(t *testing.T) {
	assert.NotPanics(t, func() {
		assertStringSliceHasAtLeast([]string{"a", "b", "c"}, 0, "Should not produce an error")
		assertStringSliceHasAtLeast([]string{"a", "b", "c"}, 1, "Should not produce an error")
		assertStringSliceHasAtLeast([]string{"a", "b", "c"}, 2, "Should not produce an error")
		assertStringSliceHasAtLeast([]string{"a", "b", "c"}, 3, "Should not produce an error")
	})
}

// TestAssertHasAtLeastWithError is a unit test.
func TestAssertHasAtLeastWithError(t *testing.T) {
	assert.Panics(t, func() {
		assertStringSliceHasAtLeast([]string{"a", "b", "c"}, 4, "Slice too small")
	})
}

// TestIsTrueIsActuallyTrue is a unit test.
func TestIsTrueIsActuallyTrue(t *testing.T) {
	assert.NotPanics(t, func() {
		assertTrue(true)
	})
}

// TestIsTrueIsActuallyFalse is a unit test.
func TestIsTrueIsActuallyFalse(t *testing.T) {
	assert.Panics(t, func() {
		assertTrue(false)
	})
}
