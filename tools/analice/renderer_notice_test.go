package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_noticeRenderer_groupByLicense(t *testing.T) {
	n := noticeRenderer{}
	m := n.groupByLicense(
		somePackage("r", true, &mit),
		somePackage("a", false, &apache2Short),
		somePackage("b", false, &bsd2),
		somePackage("n", false, nil))
	assert.NotNil(t, m)
	assert.Len(t, m, 2)
}
