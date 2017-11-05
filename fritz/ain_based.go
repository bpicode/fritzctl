package fritz

import (
	"fmt"
	"net/http"

	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
)

// ainBased API definition, guided by
// https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type ainBased interface {
	listDevices() (*Devicelist, error)
	switchOn(ain string) (string, error)
	switchOff(ain string) (string, error)
	toggle(ain string) (string, error)
	applyTemperature(value float64, ain string) (string, error)
}

// newAinBased creates a Fritz AHA API (working on AINs) from a given client.
func newAinBased(client *Client) ainBased {
	return &ainBasedClient{client: client}
}

type ainBasedClient struct {
	client *Client
}

func (a *ainBasedClient) getf(url string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		logger.Debug("GET", url)
		return a.client.HTTPClient.Get(url)
	}
}

// listDevices lists the basic data of the smart home devices.
func (a *ainBasedClient) listDevices() (*Devicelist, error) {
	url := a.homeAutoSwitch().
		query("switchcmd", "getdevicelistinfos").
		build()
	var deviceList Devicelist
	errRead := httpread.XML(a.getf(url), &deviceList)
	return &deviceList, errRead
}

// switchOn switches a device on. The device is identified by its AIN.
func (a *ainBasedClient) switchOn(ain string) (string, error) {
	return a.switchForAin(ain, "setswitchon")
}

// switchOff switches a device off. The device is identified by its AIN.
func (a *ainBasedClient) switchOff(ain string) (string, error) {
	return a.switchForAin(ain, "setswitchoff")
}

// toggle toggles the on/off state of a device. The device is identified by its AIN.
func (a *ainBasedClient) toggle(ain string) (string, error) {
	url := a.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "setswitchtoggle").
		build()
	return httpread.String(a.getf(url))
}

// applyTemperature sets the desired temperature on a "HKR" device. The device is identified by its AIN.
func (a *ainBasedClient) applyTemperature(value float64, ain string) (string, error) {
	param, err := temperatureParam(value)
	if err != nil {
		return "", err
	}
	url := a.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "sethkrtsoll").
		query("param", fmt.Sprintf("%d", param)).
		build()
	return httpread.String(a.getf(url))
}

func (a *ainBasedClient) switchForAin(ain, command string) (string, error) {
	url := a.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", command).
		build()
	return httpread.String(a.getf(url))
}

func (a *ainBasedClient) homeAutoSwitch() fritzURLBuilder {
	return newURLBuilder(a.client.Config).path(homeAutomationURI).query("sid", a.client.SessionInfo.SID)
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
