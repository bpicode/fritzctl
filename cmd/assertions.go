package cmd

import (
	"fmt"

	"github.com/pkg/errors"
)

// assertNoErr panics with log message if the argument error is not nil.
func assertNoErr(err error, format string, args ...interface{}) {
	wErr := errors.Wrapf(err, format, args...)
	assertTrue(wErr == nil, wErr)
}

// assertMinLen panics with a log message if the slice passed as argument has a size smaller than expected.
func assertMinLen(vals []string, num int, v ...interface{}) {
	assertTrue(len(vals) >= num, v...)
}

// assertTrue fails with a log message if the value is not true.
func assertTrue(val bool, v ...interface{}) {
	if !val {
		panic(fmt.Sprint(v...))
	}
}
