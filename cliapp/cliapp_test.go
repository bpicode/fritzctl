package cliapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliCreate(t *testing.T) {
	cli := Create()
	assert.NotNil(t, cli)
}
