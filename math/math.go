package math

// Round rounds a float64 value to an integer.
func Round(v float64) int64 {
	return int64(v + 0.5)
}
