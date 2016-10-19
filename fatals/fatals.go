package fatals

import "log"

// AssertNoError fails with log message if the argumetn error is not nil.
func AssertNoError(err error, v ...interface{}) {
	if err != nil {
		log.Panic(v)
	}
}
