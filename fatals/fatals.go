package fatals

import "github.com/bpicode/fritzctl/logger"

// AssertNoError fails with log message if the argument error is not nil.
func AssertNoError(err error, v ...interface{}) {
	if err != nil {
		logger.Panic(v...)
	}
}

// AssertStringSliceHasAtLeast fails with a log message if the slice passed as argument has a size smaller than expected.
func AssertStringSliceHasAtLeast(vals []string, num int, v ...interface{}) {
	if len(vals) < num {
		logger.Panic(v...)
	}
}
