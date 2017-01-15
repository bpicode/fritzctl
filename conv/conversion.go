package conv

import "strconv"

// Float64Slice wraps a float64 slice and defines additional
// operations on it.
type Float64Slice []float64

// String transforms the Float64Slice into a slice of strings
// using the specified format. See also strconv.FormatFloat.
func (fs Float64Slice) String(format byte, prec int) []string {
	var strs []string
	for _, f := range fs {
		strs = append(strs, strconv.FormatFloat(f, format, prec, 64))
	}
	return strs
}
