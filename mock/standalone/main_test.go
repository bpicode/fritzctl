package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStandalone runs the standalone main function.
func TestStandalone(t *testing.T) {
	assert.NotPanics(t, main)
}
