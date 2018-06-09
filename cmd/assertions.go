package cmd

import (
	"github.com/pkg/errors"
)

func assertNoErr(err error, format string, args ...interface{}) {
	wErr := errors.Wrapf(err, format, args...)
	assertTrue(wErr == nil, wErr)
}

func assertMinLen(vals []string, num int, format string, v ...interface{}) {
	assertTrue(len(vals) >= num, errors.Errorf(format, v...))
}

func assertTrue(val bool, err error) {
	if !val {
		panic(err)
	}
}
