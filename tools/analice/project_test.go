package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_condAppendProject(t *testing.T) {
	gmp := goModProjector{}
	m := make(map[string]dependency)

	gmp.condAppendProject(m, "")
	assert.Len(t, m, 0)

	gmp.condAppendProject(m, "github.com/fatih/color v1.7.0 h1:DkWD4oS2D8LGGgTQ6IvwJJXSL5Vp2ffcQg58nFV38Ys=")
	assert.Len(t, m, 1)
}
