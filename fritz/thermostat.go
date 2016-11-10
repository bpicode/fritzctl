package fritz

// Thermostat models the "HKR" device.
type Thermostat struct {
	Measured string `xml:"tist"`    // Measured temperature.
	Goal     string `xml:"tsoll"`   // Desired temperature, user contolled.
	Saving   string `xml:"absenk"`  // Energy saving temperature.
	Comfort  string `xml:"komfort"` // Comfortable temperature.
}
