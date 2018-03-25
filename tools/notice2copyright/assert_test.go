package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_assertTrue(t *testing.T) {
	assert.NotPanics(t, func() {
		assertTrue(log.Fatalln, true, "err")
	})
	assert.Panics(t, func() {
		assertTrue(func(v ...interface{}) {
			panic("panic")
		}, false, "err")
	})
}
