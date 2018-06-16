package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrors demonstrates usage of errors api.
func TestErrors(t *testing.T) {
	assert.Equal(t, nil, Wrapf(nil, "would-be-error"))
	err := Wrapf(fmt.Errorf("inner"), "outer")
	assert.Error(t, err)
	assert.Equal(t, "outer: inner", err.Error())

	wc, ok := err.(*withCause)
	assert.True(t, ok)
	assert.Equal(t, wc.Cause().Error(), "inner")
	assert.Equal(t, wc.Msg(), "outer")
}
