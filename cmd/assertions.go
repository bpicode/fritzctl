package cmd

import (
	"fmt"
)

// assertNoError panics with log message if the argument error is not nil.
func assertNoError(err error, v ...interface{}) {
	assertTrue(err == nil, v...)
}

// assertStringSliceHasAtLeast panics with a log message if the slice passed as argument has a size smaller than expected.
func assertStringSliceHasAtLeast(vals []string, num int, v ...interface{}) {
	assertTrue(len(vals) >= num, v...)
}

// assertTrue fails with a log message if the value is not true.
func assertTrue(val bool, v ...interface{}) {
	if !val {
		panic(fmt.Sprint(v...))
	}
}
