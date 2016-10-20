package fritz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Fritz API wrapper.
type Fritz struct {
	client *Client
}

// UsingClient is factory function to create a Fritz API interaction point.
func UsingClient(client *Client) *Fritz {
	return &Fritz{client: client}
}

func (fritz *Fritz) getWithAin(ain, switchcmd, param string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s/%s?ain=%s&switchcmd=%s&param=%s&sid=%s",
		fritz.client.Config.Protocol,
		fritz.client.Config.Host,
		"/webservices/homeautoswitch.lua",
		ain,
		switchcmd,
		param,
		fritz.client.SessionInfo.SID)
	return fritz.client.HTTPClient.Get(url)
}

func (fritz *Fritz) get(switchcmd string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s/%s?switchcmd=%s&sid=%s",
		fritz.client.Config.Protocol,
		fritz.client.Config.Host,
		"/webservices/homeautoswitch.lua",
		switchcmd,
		fritz.client.SessionInfo.SID)
	return fritz.client.HTTPClient.Get(url)
}

// GetSwitchList lists the switches configured in the system.
func (fritz *Fritz) GetSwitchList() (string, error) {
	response, errHTTP := fritz.get("getswitchlist")
	if errHTTP != nil {
		return "", errHTTP
	}
	defer response.Body.Close()
	body, errRead := ioutil.ReadAll(response.Body)
	return string(body), errRead
}

// Devicelist wraps a list of devices.
type Devicelist struct {
	Devices []Device `xml:"device"`
}

// Device models a smart home device.
type Device struct {
	Identifier      string     `xml:"identifier,attr"`
	ID              string     `xml:"id,attr"`
	Functionbitmask string     `xml:"functionbitmask,attr"`
	Fwversion       string     `xml:"fwversion,attr"`
	Manufacturer    string     `xml:"manufacturer,attr"`
	Productname     string     `xml:"productname,attr"`
	Present         int        `xml:"present"`
	Name            string     `xml:"name"`
	Switch          Switch     `xml:"switch"`
	Powermeter      Powermeter `xml:"powermeter"`
}

// Switch models the state of a switch
type Switch struct {
	State string `xml:"state"` // Switch state 1/0 on/off (empty if not known or if there was an error)
	Mode  string `xml:"mode"`  // Switch mode manual/automatic (empty if not known or if there was an error)
	Lock  string `xml:"lock"`  // Switch locked? 1/0 (empty if not known or if there was an error)
}

// Powermeter models a power measurement
type Powermeter struct {
	Power  string `xml:"power"`  // Current power, refreshed approx every 2 minutes
	Energy string `xml:"energy"` // Absolute energy consuption since the device started operating
}

// ListDevices lists the basic data of the smart home devices.
func (fritz *Fritz) ListDevices() (*Devicelist, error) {
	response, errHTTP := fritz.get("getdevicelistinfos")
	if errHTTP != nil {
		return nil, errHTTP
	}
	defer response.Body.Close()
	var deviceList Devicelist
	errDecode := xml.NewDecoder(response.Body).Decode(&deviceList)
	return &deviceList, errDecode
}
