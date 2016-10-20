package fatals

import "github.com/bpicode/fritzctl/logger"

// AssertNoError fails with log message if the argumetn error is not nil.
func AssertNoError(err error, v ...interface{}) {
	if err != nil {
		logger.Panic(v...)
	}
}
