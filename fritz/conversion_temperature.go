package fritz

import "strconv"

func fmtTemperatureHkr(th string) string {
	f, err := strconv.ParseFloat(th, 64)
	if err != nil {
		return ""
	}
	return fmtTemperatureWithSpecialBoundaries(f)
}

func fmtTemperatureWithSpecialBoundaries(f float64) string {
	switch {
	case f == 255:
		return "?"
	case f == 254:
		return "ON"
	case f == 253:
		return "OFF"
	case f < 16:
		return "8"
	case f > 56:
		return "28"
	default:
		return strconv.FormatFloat(f*0.5, 'f', -1, 64)
	}
}
