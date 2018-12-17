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

// Keys extracts the key of a map[string]string and returns them
// as a slice of strings.
func Keys(m map[string]string) []string {
	return Contract(m, func(k string, v string) string {
		return k
	})
}

// Contract takes a map[string]string and contracts pairs of key and values.
func Contract(m map[string]string, f func(string, string) string) []string {
	c := make([]string, 0, len(m))
	for k, v := range m {
		c = append(c, f(k, v))
	}
	return c
}

// ErrorMessages accumulates the error messages from slice of errors.
func ErrorMessages(errs []error) []string {
	ms := make([]string, 0, len(errs))
	for _, err := range errs {
		ms = append(ms, err.Error())
	}
	return ms
}

// MapWithDefault provides a closure that does the regular map lookup, and defaults to the passed value no value is found.
func MapWithDefault(m map[string]string, def string) func(string) string {
	return func(key string) string {
		v, ok := m[key]
		if !ok {
			return def
		}
		return v
	}
}
