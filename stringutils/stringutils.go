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

// StringKeys extracts the key of a map[string]string and returns them
// as a slice of strings.
func StringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// StringValues extracts the values of a map[string]string and returns them
// as a slice of strings.
func StringValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

// Filter takes a slice of strings and returns a slice of strings
// consisting of those strings from the input for which the filte
// evaluated to true.
func Filter(strs []string, f func(string) bool) []string {
	filtered := make([]string, 0, len(strs))
	for _, str := range strs {
		if f(str) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

// ErrorMessages accumulates the error messages from s lice of errors.
func ErrorMessages(errs []error) []string {
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return msgs
}
