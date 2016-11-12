package fritz

// Temperature models a temperature measurement.
type Temperature struct {
	Celsius string `xml:"celsius"` // Temperature measured at the device sensor in units of 0.1 °C. Negative and positive values are possible.
	Offset  string `xml:"offset"`  // Temperature offset (set by the user) in units of 0.1 °C. Negative and positive values are possible.
}
