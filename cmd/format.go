package cmd

import (
	"fmt"
	"strconv"
)

func fmtUnit(f func() string, unit string) string {
	s := f()
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%s %s", s, unit)
}
