package fritz

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/httpread"
)

// Fritz API definition.
type Fritz interface {
	GetSwitchList() (string, error)
	ListDevices() (*Devicelist, error)
	GetAinForName(name string) (string, error)
	Switch(name, state string) (string, error)
	Toggle(name string) (string, error)
	Temperature(name string, value float64) (string, error)
}

// fritzImpl implements Fritz API.
type fritzImpl struct {
	client *Client
}

// UsingClient is factory function to create a Fritz API interaction point.
func UsingClient(client *Client) Fritz {
	return &fritzImpl{client: client}
}

func (fritz *fritzImpl) getWithAinAndParam(ain, switchcmd, param string) (*http.Response, error) {
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

func (fritz *fritzImpl) getWithAin(ain, switchcmd string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s/%s?ain=%s&switchcmd=%s&sid=%s",
		fritz.client.Config.Protocol,
		fritz.client.Config.Host,
		"/webservices/homeautoswitch.lua",
		ain,
		switchcmd,
		fritz.client.SessionInfo.SID)
	return fritz.client.HTTPClient.Get(url)
}

func (fritz *fritzImpl) get(switchcmd string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s/%s?switchcmd=%s&sid=%s",
		fritz.client.Config.Protocol,
		fritz.client.Config.Host,
		"/webservices/homeautoswitch.lua",
		switchcmd,
		fritz.client.SessionInfo.SID)
	return fritz.client.HTTPClient.Get(url)
}

// GetSwitchList lists the switches configured in the system.
func (fritz *fritzImpl) GetSwitchList() (string, error) {
	response, errHTTP := fritz.get("getswitchlist")
	return httpread.ReadFullyString(response, errHTTP)
}

// Devicelist wraps a list of devices.
type Devicelist struct {
	Devices []Device `xml:"device"`
}

// Device models a smart home device.
type Device struct {
	Identifier      string      `xml:"identifier,attr"`
	ID              string      `xml:"id,attr"`
	Functionbitmask string      `xml:"functionbitmask,attr"`
	Fwversion       string      `xml:"fwversion,attr"`
	Manufacturer    string      `xml:"manufacturer,attr"`
	Productname     string      `xml:"productname,attr"`
	Present         int         `xml:"present"`
	Name            string      `xml:"name"`
	Switch          Switch      `xml:"switch"`
	Powermeter      Powermeter  `xml:"powermeter"`
	Temperature     Temperature `xml:"temperature"`
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

// Temperature models a temperature measurement.
type Temperature struct {
	Celsius string `xml:"celsius"` // Current power, refreshed approx every 2 minutes
	Offset  string `xml:"offset"`  // Absolute energy consuption since the device started operating
}

// ListDevices lists the basic data of the smart home devices.
func (fritz *fritzImpl) ListDevices() (*Devicelist, error) {
	response, errHTTP := fritz.get("getdevicelistinfos")
	if errHTTP != nil {
		return nil, errHTTP
	}
	defer response.Body.Close()
	var deviceList Devicelist
	errDecode := xml.NewDecoder(response.Body).Decode(&deviceList)
	return &deviceList, errDecode
}

// Switch turns a device on/off.
func (fritz *fritzImpl) Switch(name, state string) (string, error) {
	ain, errGetAin := fritz.GetAinForName(name)
	if errGetAin != nil {
		return "", errGetAin
	}
	return fritz.switchForAin(ain, state)
}

func (fritz *fritzImpl) switchForAin(ain, state string) (string, error) {
	resp, errSwitch := fritz.getWithAin(ain, switchCommandFor(state))
	return httpread.ReadFullyString(resp, errSwitch)
}

// GetAinForName returns the AIN corresponding to a device name.
func (fritz *fritzImpl) GetAinForName(name string) (string, error) {
	devList, errList := fritz.ListDevices()
	if errList != nil {
		return "", errList
	}
	devs := devList.Devices
	names := make([]string, len(devs))
	for i, dev := range devs {
		names[i] = dev.Name
	}

	var ain string
	for _, dev := range devs {
		if dev.Name == name {
			ain = strings.Replace(dev.Identifier, " ", "", -1)
		}
	}
	if ain == "" {
		return "", errors.New("No device found with name '" + name + "'. Available devices are " + fmt.Sprintf("%s", names))
	}
	return ain, nil
}

func switchCommandFor(state string) string {
	if strings.EqualFold(state, "on") {
		return "setswitchon"
	}
	return "setswitchoff"
}

// Toggle toggles the on/off state of a device.
func (fritz *fritzImpl) Toggle(name string) (string, error) {
	ain, errGetAin := fritz.GetAinForName(name)
	if errGetAin != nil {
		return "", errGetAin
	}
	return fritz.toggleForAin(ain)
}

func (fritz *fritzImpl) toggleForAin(ain string) (string, error) {
	resp, errSwitch := fritz.getWithAin(ain, "setswitchtoggle")
	return httpread.ReadFullyString(resp, errSwitch)
}

// Temperature sets the desired temperature of a "HKR" device.
func (fritz *fritzImpl) Temperature(name string, value float64) (string, error) {
	ain, errGetAin := fritz.GetAinForName(name)
	if errGetAin != nil {
		return "", errGetAin
	}
	return fritz.temperatureForAin(ain, value)
}

func (fritz *fritzImpl) temperatureForAin(ain string, value float64) (string, error) {
	doubledValue := 2 * value
	rounded := round(doubledValue)
	response, err := fritz.getWithAinAndParam(ain, "sethkrtsoll", fmt.Sprintf("%d", rounded))
	return httpread.ReadFullyString(response, err)
}

func round(v float64) int64 {
	return int64(v + 0.5)
}
