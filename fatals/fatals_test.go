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
