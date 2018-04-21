package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseCopyrightHolder(t *testing.T) {
	e := copyrightEngine{}
	assert.Equal(t, "2016 bpicode", e.parseCopyrightHolder("Copyright (c) 2016 bpicode"))
	assert.Equal(t, "John Doe <john.d@example.com>", e.parseCopyrightHolder("Copyright (c) John Doe <john.d@example.com>"))
	assert.Equal(t, "2009 The Go Authors. All rights reserved.", e.parseCopyrightHolder("Copyright (c) 2009 The Go Authors. All rights reserved."))
	assert.Equal(t, "2011 John Doe <john.d@example.com>", e.parseCopyrightHolder("Copyright © 2011 John Doe <john.d@example.com>"))
}
