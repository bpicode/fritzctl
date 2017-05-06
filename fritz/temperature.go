package fritz

import "strconv"

// Temperature models a temperature measurement.
type Temperature struct {
	Celsius string `xml:"celsius"` // Temperature measured at the device sensor in units of 0.1 °C. Negative and positive values are possible.
	Offset  string `xml:"offset"`  // Temperature offset (set by the user) in units of 0.1 °C. Negative and positive values are possible.
}

// FmtCelsius formats the value of t.Celsius as obtained on the http interface as a stringified floating point number.
func (t *Temperature) FmtCelsius() string {
	return t.fmtTemperature(t.Celsius)
}

// FmtOffset formats the value of t.Offset as obtained on the http interface as a stringified floating point number.
func (t *Temperature) FmtOffset() string {
	return t.fmtTemperature(t.Offset)
}

func (t *Temperature) fmtTemperature(th string) string {
	f, err := strconv.ParseFloat(th, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(f*0.1, 'f', -1, 64)
}
