package fritz

// LanDevices is the top-level wrapper for the FRITZ!Box answer upon
// a query for lan devices.
type LanDevices struct {
	Network []NetworkElement `json:"network"`
}

// NetworkElement corresponds to a single entry in LanDevices.
type NetworkElement struct {
	Name                 string `json:"name"`
	IP                   string `json:"ip"`
	Mac                  string `json:"mac"`
	UID                  string `json:"UID"`
	Dhcp                 string `json:"dhcp"`
	Wlan                 string `json:"wlan"`
	Ethernet             string `json:"ethernet"`
	Active               string `json:"active"`
	StaticDhcp           string `json:"static_dhcp"`
	ManuName             string `json:"manu_name"`
	Wakeup               string `json:"wakeup"`
	Deleteable           string `json:"deleteable"`
	Source               string `json:"source"`
	Online               string `json:"online"`
	Speed                string `json:"speed"`
	WlanUIDs             string `json:"wlan_UIDs"`
	AutoWakeup           string `json:"auto_wakeup"`
	Guest                string `json:"guest"`
	URL                  string `json:"url"`
	WlanStationType      string `json:"wlan_station_type"`
	EthernetPort         string `json:"ethernet_port"`
	WlanShowInMonitor    string `json:"wlan_show_in_monitor"`
	Plc                  string `json:"plc"`
	ParentalControlAbuse string `json:"parental_control_abuse"`
}
