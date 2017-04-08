package fritz

// Switch models the state of a switch.
type Switch struct {
	State      string `xml:"state"`      // Switch state 1/0 on/off (empty if not known or if there was an error).
	Mode       string `xml:"mode"`       // Switch mode manual/automatic (empty if not known or if there was an error).
	Lock       string `xml:"lock"`       // Switch locked (box defined)? 1/0 (empty if not known or if there was an error).
	DeviceLock string `xml:"devicelock"` // Switch locked (device defined)? 1/0 (empty if not known or if there was an error).
}
