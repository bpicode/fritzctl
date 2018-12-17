package fritz

// AlertSensor collects data from devices that have an alert sensor.
type AlertSensor struct {
	State string `xml:"state"` // Last transmitted alert state, "0" - no alert, "1" - alert, "" if unknown or upon errors.
}
