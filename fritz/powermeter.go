package fritz

// Powermeter models a power measurement
type Powermeter struct {
	Power  string `xml:"power"`  // Current power, refreshed approx every 2 minutes
	Energy string `xml:"energy"` // Absolute energy consuption since the device started operating
}
