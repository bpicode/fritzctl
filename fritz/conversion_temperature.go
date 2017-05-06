package fritz

import "strconv"

func fmtTemperatureHkr(th string) string {
	f, err := strconv.ParseFloat(th, 64)
	if err != nil {
		return ""
	}
	var str string
	switch {
	case f == 255:
		str = "?"
	case f == 254:
		str = "ON"
	case f == 253:
		str = "OFF"
	case f < 16:
		str = "8"
	case f > 56:
		str = "28"
	default:
		str = strconv.FormatFloat(f*0.5, 'f', -1, 64)
	}
	return str
}
