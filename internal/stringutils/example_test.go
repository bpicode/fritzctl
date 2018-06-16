package stringutils_test

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bpicode/fritzctl/internal/stringutils"
)

// Transform applies a given operation to all elements of a string slice.
func ExampleTransform() {
	s := []string{"hello", "world"}
	u := stringutils.Transform(s, strings.ToUpper)
	fmt.Println(u)
	// output: [HELLO WORLD]
}

// Quote puts double quotation marks around each element of a string slice.
func ExampleQuote() {
	s := []string{"hello", "world"}
	q := stringutils.Quote(s)
	fmt.Println(q)
	// output: ["hello" "world"]
}

// DefaultIf takes three strings v, d, c and returns (v != c)? v : d.
func ExampleDefaultIf() {
	s := stringutils.DefaultIf("xxx", "[CENSORED]", "xxx")
	fmt.Println(s)
	s = stringutils.DefaultIf("abc", "[CENSORED]", "xxx")
	fmt.Println(s)
	// output: [CENSORED]
	// abc
}

// DefaultIfEmpty takes two strings v, d and returns (v != "")? v : d.
func ExampleDefaultIfEmpty() {
	s := stringutils.DefaultIfEmpty("my string", "<no value>")
	fmt.Println(s)
	s = stringutils.DefaultIfEmpty("", "<no value>")
	fmt.Println(s)
	// output: my string
	// <no value>
}

// Keys extracts the keys of a map[string]string.
func ExampleKeys() {
	m := map[string]string{"key1": "value1", "key2": "value2"}
	ks := stringutils.Keys(m)
	sort.Strings(ks)
	fmt.Println(ks)
	// output: [key1 key2]
}

// ErrorMessages extracts the messages from a slice of errors.
func ExampleErrorMessages() {
	es := []error{errors.New("an_error"), errors.New("another_error")}
	ms := stringutils.ErrorMessages(es)
	fmt.Println(ms)
	// output: [an_error another_error]
}

// Contract takes a map[string]string and pipes key-value pairs through a function and returns accumulated results.
func ExampleContract() {
	m := map[string]string{"key1": "value1", "key2": "value2"}
	vs := stringutils.Contract(m, func(k string, v string) string { return v })
	sort.Strings(vs)
	fmt.Println(vs)
	// output: [value1 value2]
}
