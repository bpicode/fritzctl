package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadStdlib(t *testing.T) {
	s := loadStdlib()
	assert.NotNil(t, s)
	assert.True(t, len(s) > 0)
}
