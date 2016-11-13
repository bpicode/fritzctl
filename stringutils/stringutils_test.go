package stringutils

import (
	"fmt"
	"strings"
	"testing"

	"errors"

	"sort"

	"github.com/stretchr/testify/assert"
)

// TestTransform tests the transformation by transforming strings to uppercase.
func TestTransform(t *testing.T) {
	strs := []string{"a", "b", "c"}
	allupper := Transform(strs, strings.ToUpper)
	assert.Len(t, allupper, len(strs))
	assert.Equal(t, "A", allupper[0])
	assert.Equal(t, "B", allupper[1])
	assert.Equal(t, "C", allupper[2])
}

// TestQuote tests quoting of strings.
func TestQuote(t *testing.T) {
	strs := []string{"a", "b", "c"}
	quoted := Quote(strs)
	assert.Len(t, quoted, len(strs))
	assert.Equal(t, "\"a\"", quoted[0])
	assert.Equal(t, "\"b\"", quoted[1])
	assert.Equal(t, "\"c\"", quoted[2])
}

// TestDefaultIfEmpty tests several constallations of default value returns.
func TestDefaultIfEmpty(t *testing.T) {
	assert.Equal(t, "AA", DefaultIfEmpty("AA", "BB"))
	assert.Equal(t, "BB", DefaultIfEmpty("", "BB"))
	assert.Equal(t, "", DefaultIfEmpty("", ""))
	var a string
	assert.Equal(t, "XX", DefaultIfEmpty(a, "XX"))
}

// TestStringKeysAndValues tests string key/value extraction on several examples.
func TestStringKeysAndValues(t *testing.T) {
	testCases := []struct {
		input          map[string]string
		expectedKeys   []string
		expectedValues []string
	}{
		{input: map[string]string{}, expectedKeys: []string{}, expectedValues: []string{}},
		{input: map[string]string{"foo": "a", "bar": "b"}, expectedKeys: []string{"foo", "bar"}, expectedValues: []string{"a", "b"}},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test string key extraction %d", i), func(t *testing.T) {
			keys := StringKeys(testCase.input)
			assert.NotNil(t, keys)
			sort.Strings(testCase.expectedKeys)
			sort.Strings(keys)
			assert.Equal(t, testCase.expectedKeys, keys)

			values := StringValues(testCase.input)
			assert.NotNil(t, values)
			sort.Strings(testCase.expectedValues)
			sort.Strings(values)
			assert.Equal(t, testCase.expectedValues, values)
		})
	}
}

// TestFilter tests the filter function on various examples.
func TestFilter(t *testing.T) {
	testCases := []struct {
		input          []string
		filter         func(string) bool
		expectedOutput []string
	}{
		{input: []string{}, filter: func(string) bool { return true }, expectedOutput: []string{}},
		{input: []string{"a", "b", "c"}, filter: func(string) bool { return true }, expectedOutput: []string{"a", "b", "c"}},
		{input: []string{"a", "b", "c"}, filter: func(string) bool { return false }, expectedOutput: []string{}},
		{input: []string{"a", "b", "c"}, filter: func(s string) bool { return s == "b" }, expectedOutput: []string{"b"}},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test string filter %d", i), func(t *testing.T) {
			output := Filter(testCase.input, testCase.filter)
			assert.NotNil(t, output)
			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}

// TestErrorMessages tests the error message extraction  on various examples.
func TestErrorMessages(t *testing.T) {
	testCases := []struct {
		input          []error
		expectedOutput []string
	}{
		{input: []error{}, expectedOutput: []string{}},
		{input: []error{errors.New("some error")}, expectedOutput: []string{"some error"}},
		{input: []error{errors.New("some error"), errors.New("another error")}, expectedOutput: []string{"some error", "another error"}},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test errpr message extraction %d", i), func(t *testing.T) {
			output := ErrorMessages(testCase.input)
			assert.NotNil(t, output)
			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}
