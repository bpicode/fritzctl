package fatals

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertNoErrorWithNoError(t *testing.T) {
	AssertNoError(nil)
}

func TestAssertNoErrorWithError(t *testing.T) {
	assert.Panics(t, func() {
		AssertNoError(errors.New("sum ting wong"))
	})
}
