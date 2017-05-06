package fritz

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/fritzclient"
	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
)

// HomeAutomationApi API definition, guided by
// https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type HomeAutomationApi interface {
	ListDevices() (*Devicelist, error)
	NameToAinTable() (map[string]string, error)
	SwitchOn(ain string) (string, error)
	SwitchOff(ain string) (string, error)
	Toggle(ain string) (string, error)
	ApplyTemperature(value float64, ain string) (string, error)
}

// HomeAutomation creates a Fritz AHA API from a given client.
func HomeAutomation(client *fritzclient.Client) HomeAutomationApi {
	return &ahaHttp{client: client}
}

type ahaHttp struct {
	client *fritzclient.Client
}

func (aha *ahaHttp) getf(url string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		logger.Debug("GET", url)
		return aha.client.HTTPClient.Get(url)
	}
}

// ListDevices lists the basic data of the smart home devices.
func (aha *ahaHttp) ListDevices() (*Devicelist, error) {
	url := aha.homeAutoSwitch().
		query("switchcmd", "getdevicelistinfos").
		build()
	var deviceList Devicelist
	errRead := httpread.ReadFullyXML(aha.getf(url), &deviceList)
	return &deviceList, errRead
}

// SwitchOn switches a device on. The device is identified by its AIN.
func (aha *ahaHttp) SwitchOn(ain string) (string, error) {
	return aha.switchForAin(ain, "setswitchon")
}

// SwitchOff switches a device off. The device is identified by its AIN.
func (aha *ahaHttp) SwitchOff(ain string) (string, error) {
	return aha.switchForAin(ain, "setswitchoff")
}

// Toggle toggles the on/off state of a device. The device is identified by its AIN.
func (aha *ahaHttp) Toggle(ain string) (string, error) {
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "setswitchtoggle").
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

// ApplyTemperature sets the desired temperature on a "HKR" device. The device is identified by its AIN.
func (aha *ahaHttp) ApplyTemperature(value float64, ain string) (string, error) {
	doubledValue := 2 * value
	rounded := math.Round(doubledValue)
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "sethkrtsoll").
		query("param", fmt.Sprintf("%d", rounded)).
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

func (aha *ahaHttp) switchForAin(ain, command string) (string, error) {
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", command).
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

// NameToAinTable returns a lookup name -> AIN.
func (aha *ahaHttp) NameToAinTable() (map[string]string, error) {
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

func (aha *ahaHttp) homeAutoSwitch() fritzURLBuilder {
	return newURLBuilder(aha.client.Config).path(homeAutomationURI).query("sid", aha.client.SessionInfo.SID)
}
