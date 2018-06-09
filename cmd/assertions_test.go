package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAssertNoErrorPassesOnNil tests assertNoErr, "happy path".
func TestAssertNoErrorPassesOnNil(t *testing.T) {
	assert.NotPanics(t, func() {
		assertNoErr(nil, "would-be-context")
	})
}

// TestAssertNoErrorPanicsOnError tests assertNoErr, "error path".
func TestAssertNoErrorPanicsOnError(t *testing.T) {
	assert.Panics(t, func() {
		assertNoErr(errors.New("we too low"), "sum ting wong")
	})
}

// TestAssertMinLenPassesOnVerification tests assertMinLen, "happy path".
func TestAssertMinLenPassesOnVerification(t *testing.T) {
	assert.NotPanics(t, func() {
		assertMinLen([]string{"a", "b", "c"}, 0, "Should not produce an error")
		assertMinLen([]string{"a", "b", "c"}, 1, "Should not produce an error")
		assertMinLen([]string{"a", "b", "c"}, 2, "Should not produce an error")
		assertMinLen([]string{"a", "b", "c"}, 3, "Should not produce an error")
	})
}

// TestAssertMinLenPanicsOnViolation tests assertMinLen, "error path".
func TestAssertMinLenPanicsOnViolation(t *testing.T) {
	assert.Panics(t, func() {
		assertMinLen([]string{"a", "b", "c"}, 4, "Slice too small")
	})
}

// TestIsTruePassesOnVerification tests assertTrue, "happy path".
func TestIsTruePassesOnVerification(t *testing.T) {
	assert.NotPanics(t, func() {
		assertTrue(true, errors.New("would-be-error"))
	})
}

// TestIsTruePanicsOnViolation tests assertTrue, "error path".
func TestIsTruePanicsOnViolation(t *testing.T) {
	assert.Panics(t, func() {
		assertTrue(false, errors.New("is-an-error"))
	})
}
