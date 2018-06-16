package cmd

import (
	"fmt"

	"github.com/bpicode/fritzctl/internal/errors"
)

func assertNoErr(err error, format string, args ...interface{}) {
	wErr := errors.Wrapf(err, format, args...)
	assertTrue(wErr == nil, wErr)
}

func assertMinLen(vals []string, num int, format string, v ...interface{}) {
	assertTrue(len(vals) >= num, fmt.Errorf(format, v...))
}

func assertTrue(val bool, err error) {
	if !val {
		panic(err)
	}
}
