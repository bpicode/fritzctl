package stringutils

// Transform takes a slice of strings, transforms each element
// and returns a slice of strings containing the transformed values
func Transform(s []string, f func(string) string) []string {
	var transformed []string
	for _, e := range s {
		transformed = append(transformed, f(e))
	}
	return transformed
}

// Quote takes a slice of strings and returns a slice
// of strings where each element of the input is put
// in double quotation marks.
func Quote(s []string) []string {
	return Transform(s, func(str string) string {
		return `"` + str + `"`
	})
}

// DefaultIfEmpty falls back to a default value if the passed value is empty
func DefaultIfEmpty(value, defaultValue string) string {
	return DefaultIf(value, defaultValue, "")
}

// DefaultIf falls back to a default value if the passed value is equal
// to the condition parameter.
func DefaultIf(value, defaultValue, condition string) string {
	if value == condition {
		return defaultValue
	}
	return value
}
