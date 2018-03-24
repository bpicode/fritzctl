package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_openFirst(t *testing.T) {
	_, err := openFirst("/that/should/not/work/ever", "asdfagag.txt", "asasgdsg.md")
	assert.Error(t, err)
}

func Test_parseCopyrightHolder(t *testing.T) {
	assert.Equal(t, "2016 bpicode", parseCopyrightHolder("Copyright (c) 2016 bpicode"))
	assert.Equal(t, "John Doe <john.d@example.com>", parseCopyrightHolder("Copyright (c) John Doe <john.d@example.com>"))
	assert.Equal(t, "2009 The Go Authors. All rights reserved.", parseCopyrightHolder("Copyright (c) 2009 The Go Authors. All rights reserved."))
	assert.Equal(t, "2011 John Doe <john.d@example.com>", parseCopyrightHolder("Copyright © 2011 John Doe <john.d@example.com>"))
}
