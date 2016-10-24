package units

import "strconv"

// ParseFloatAndScale parses a string to a float64, applies a scaling factor to it and returns the formatted result.
func ParseFloatAndScale(str string, scale float64) string {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(f*scale, 'f', -1, 64)
}
