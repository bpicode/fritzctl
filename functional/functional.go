package functional

import (
	"errors"
	"fmt"
)

// FirstWithoutError tries the passed functions in order and returns the first value
// obtained without error. If all functions have an error, FirstWithoutError returns
// itself an emcompassing error.
func FirstWithoutError(fcs ...func() (string, error)) (string, error) {
	var ret string
	var err error
	var errs []error
	for _, f := range fcs {
		ret, err = f()
		if err == nil {
			return ret, nil
		}
		errs = append(errs, err)
	}
	return "", errors.New(fmt.Sprint(errs))
}

// Curry a single-arg string -> (string, error) function. The returned value is
// effectively a function with no arguments.
func Curry(arg string, f func(string) (string, error)) func() (string, error) {
	return func() (string, error) {
		return f(arg)
	}
}

// Compose multiple (string) -> (string, error) functions.
func Compose(arg0 string, fcs ...func(string) (string, error)) func() (string, error) {
	return func() (string, error) {
		arg := arg0
		for _, f := range fcs {
			var err error
			arg, err = f(arg)
			if err != nil {
				return arg, err
			}
		}
		return arg, nil
	}
}
