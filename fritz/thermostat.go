package fritz

// Thermostat models the "HKR" device.
// codebeat:disable[TOO_MANY_IVARS]
type Thermostat struct {
	Measured   string `xml:"tist"`           // Measured temperature.
	Goal       string `xml:"tsoll"`          // Desired temperature, user controlled.
	Saving     string `xml:"absenk"`         // Energy saving temperature.
	Comfort    string `xml:"komfort"`        // Comfortable temperature.
	NextChange NextChange `xml:"nextchange"` // The next scheduled temperature change.
	Lock       string `xml:"lock"`           // Switch locked (box defined)? 1/0 (empty if not known or if there was an error).
	DeviceLock string `xml:"devicelock"`     // Switch locked (device defined)? 1/0 (empty if not known or if there was an error).
	ErrorCode  string `xml:"errorcode"`      // Error codes: 0 = OK, 1 = ... see https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
	BatteryLow string `xml:"batterylow"`     // "0" if the battery is OK, "1" if it is running low on capacity.
}

type NextChange struct {
	TimeStamp string `xml:"endperiod"` // Timestamp (epoch time) when the next temperature switch is scheduled.
	Goal      string `xml:"tchange"`   // The temperature to switch to. Same unit convention as in Thermostat.Measured.
}

// codebeat:enable[TOO_MANY_IVARS]
