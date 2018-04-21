package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_noticeCmd_no_project(t *testing.T) {
	err := genNotice(noticeCmd, []string{"/points/to/nowhere"})
	assert.Error(t, err)
}
