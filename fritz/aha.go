package fritz

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
)

// HomeAutomationAPI API definition, guided by
// https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type HomeAutomationAPI interface {
	ListDevices() (*Devicelist, error)
	NameToAinTable() (map[string]string, error)
	SwitchOn(ain string) (string, error)
	SwitchOff(ain string) (string, error)
	Toggle(ain string) (string, error)
	ApplyTemperature(value float64, ain string) (string, error)
}

// HomeAutomation creates a Fritz AHA API from a given client.
func HomeAutomation(client *Client) HomeAutomationAPI {
	return &ahaHTTP{client: client}
}

type ahaHTTP struct {
	client *Client
}

func (aha *ahaHTTP) getf(url string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		logger.Debug("GET", url)
		return aha.client.HTTPClient.Get(url)
	}
}

// ListDevices lists the basic data of the smart home devices.
func (aha *ahaHTTP) ListDevices() (*Devicelist, error) {
	url := aha.homeAutoSwitch().
		query("switchcmd", "getdevicelistinfos").
		build()
	var deviceList Devicelist
	errRead := httpread.ReadFullyXML(aha.getf(url), &deviceList)
	return &deviceList, errRead
}

// SwitchOn switches a device on. The device is identified by its AIN.
func (aha *ahaHTTP) SwitchOn(ain string) (string, error) {
	return aha.switchForAin(ain, "setswitchon")
}

// SwitchOff switches a device off. The device is identified by its AIN.
func (aha *ahaHTTP) SwitchOff(ain string) (string, error) {
	return aha.switchForAin(ain, "setswitchoff")
}

// Toggle toggles the on/off state of a device. The device is identified by its AIN.
func (aha *ahaHTTP) Toggle(ain string) (string, error) {
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "setswitchtoggle").
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

// ApplyTemperature sets the desired temperature on a "HKR" device. The device is identified by its AIN.
func (aha *ahaHTTP) ApplyTemperature(value float64, ain string) (string, error) {
	param, err := temperatureParam(value)
	if err != nil {
		return "", err
	}
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "sethkrtsoll").
		query("param", fmt.Sprintf("%d", param)).
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

func (aha *ahaHTTP) switchForAin(ain, command string) (string, error) {
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", command).
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

// NameToAinTable returns a lookup name -> AIN.
func (aha *ahaHTTP) NameToAinTable() (map[string]string, error) {
	devList, err := aha.ListDevices()
	if err != nil {
		return nil, err
	}
	devs := devList.Devices
	table := make(map[string]string, len(devs))
	for _, dev := range devs {
		table[dev.Name] = strings.Replace(dev.Identifier, " ", "", -1)
	}
	return table, nil
}

func (aha *ahaHTTP) homeAutoSwitch() fritzURLBuilder {
	return newURLBuilder(aha.client.Config).path(homeAutomationURI).query("sid", aha.client.SessionInfo.SID)
}

func temperatureParam(t float64) (int64, error) {
	doubled := round(2 * t)
	regular := doubled >= 16 && doubled <= 56
	special := doubled == 253 || doubled == 254
	if !(regular || special) {
		return 0, fmt.Errorf("invalid temperature value: %.1f°C is not contained in the set of acceptable values: 8-28°C, 126.5, 127", t)
	}
	return doubled, nil
}

// round rounds a float64 value to an integer.
func round(v float64) int64 {
	return int64(v + 0.5)
}
