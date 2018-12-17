package stringutils

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

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
		input        map[string]string
		expectedKeys []string
	}{
		{input: map[string]string{}, expectedKeys: []string{}},
		{input: map[string]string{"foo": "a", "bar": "b"}, expectedKeys: []string{"foo", "bar"}},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test string key extraction %d", i), func(t *testing.T) {
			keys := Keys(testCase.input)
			assert.NotNil(t, keys)
			sort.Strings(testCase.expectedKeys)
			sort.Strings(keys)
			assert.Equal(t, testCase.expectedKeys, keys)
		})
	}
}

// TestErrorMessages tests the error message extraction on various examples.
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

// TestContract tests the string contraction on various examples.
func TestContract(t *testing.T) {
	testCases := []struct {
		m      map[string]string
		f      func(string, string) string
		expect []string
	}{
		{m: map[string]string{}, f: nil, expect: []string{}},
		{m: map[string]string{"k": "v"}, f: func(a, b string) string { return a }, expect: []string{"k"}},
		{m: map[string]string{"k": "v"}, f: func(a, b string) string { return b }, expect: []string{"v"}},
		{m: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}, f: func(a, b string) string { return a + "=" + b }, expect: []string{"k1=v1", "k2=v2", "k3=v3"}},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test string contraction %d", i), func(t *testing.T) {
			out := Contract(testCase.m, testCase.f)
			assert.NotNil(t, out)
			for _, e := range testCase.expect {
				assert.Contains(t, out, e)
			}
		})
	}
}

// TestMapWithDefault probes the default valued map.
func TestMapWithDefault(t *testing.T) {
	f := MapWithDefault(map[string]string{"a": "b", "foo": "bar"}, "x")
	assert.Equal(t, "b", f("a"))
	assert.Equal(t, "bar", f("foo"))
	assert.Equal(t, "x", f("everything else"))
}
