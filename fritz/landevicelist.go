package fritz

// LanDevices is the top-level wrapper for the FRITZ!Box answer upon
// a query for lan devices.
type LanDevices struct {
	Network []NetworkElement `json:"network"`
}

// NetworkElement corresponds to a single entry in LanDevices.
// codebeat:disable[TOO_MANY_IVARS]
type NetworkElement struct {
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Mac        string `json:"mac"`
	UID        string `json:"UID"`
	Dhcp       string `json:"dhcp"`
	Wlan       string `json:"wlan"`
	Ethernet   string `json:"ethernet"`
	Active     string `json:"active"`
	Wakeup     string `json:"wakeup"`
	Deleteable string `json:"deleteable"`
	Source     string `json:"source"`
	Online     string `json:"online"`
	Speed      string `json:"speed"`
	Guest      string `json:"guest"`
	URL        string `json:"url"`
}

// codebeat:enable[TOO_MANY_IVARS]
