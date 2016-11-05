package math

import "strconv"

// ParseFloatAndScale parses a string to a float64, applies a scaling factor to it and returns the formatted result.
func ParseFloatAndScale(str string, scale float64) string {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(f*scale, 'f', -1, 64)
}

// ParseFloatAddAndScale parses two strings, adds them, applies a scaling factor to it and returns the formatted result.
func ParseFloatAddAndScale(one, another string, scale float64) string {
	first, err := strconv.ParseFloat(one, 64)
	if err != nil {
		return ""
	}
	second, err := strconv.ParseFloat(another, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat((first+second)*scale, 'f', -1, 64)
}

// Round rounds a float64 value to an integer.
func Round(v float64) int64 {
	return int64(v + 0.5)
}
