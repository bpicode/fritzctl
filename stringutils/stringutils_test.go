package stringutils

import "testing"
import "strings"
import "github.com/stretchr/testify/assert"

// TestTransform is a unit test.
func TestTransform(t *testing.T) {
	strs := []string{"a", "b", "c"}
	allupper := Transform(strs, strings.ToUpper)
	assert.Len(t, allupper, len(strs))
	assert.Equal(t, "A", allupper[0])
	assert.Equal(t, "B", allupper[1])
	assert.Equal(t, "C", allupper[2])
}

// TestQuote is a unit test.
func TestQuote(t *testing.T) {
	strs := []string{"a", "b", "c"}
	quoted := Quote(strs)
	assert.Len(t, quoted, len(strs))
	assert.Equal(t, "\"a\"", quoted[0])
	assert.Equal(t, "\"b\"", quoted[1])
	assert.Equal(t, "\"c\"", quoted[2])
}

// TestDefaultValue is a unit test.
func TestDefaultValue(t *testing.T) {
	assert.Equal(t, "AA", DefaultIfEmpty("AA", "BB"))
	assert.Equal(t, "BB", DefaultIfEmpty("", "BB"))
	assert.Equal(t, "", DefaultIfEmpty("", ""))
	var a string
	assert.Equal(t, "XX", DefaultIfEmpty(a, "XX"))
}
