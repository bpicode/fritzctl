package assert

import "github.com/bpicode/fritzctl/logger"

// NoError panics with log message if the argument error is not nil.
func NoError(err error, v ...interface{}) {
	IsTrue(err == nil, v...)
}

// StringSliceHasAtLeast panics with a log message if the slice passed as argument has a size smaller than expected.
func StringSliceHasAtLeast(vals []string, num int, v ...interface{}) {
	IsTrue(len(vals) >= num, v...)
}

// IsTrue fails with a log message if the value is not true.
func IsTrue(val bool, v ...interface{}) {
	if !val {
		logger.Panic(v...)
	}
}
