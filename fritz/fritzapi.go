package fritz

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"sync/atomic"

	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
	"github.com/bpicode/fritzctl/stringutils"
)

// Fritz API definition, guided by
// https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type Fritz interface {
	ListDevices() (*Devicelist, error)
	GetAinForName(name string) (string, error)
	SwitchOn(name string) (string, error)
	SwitchOff(name string) (string, error)
	Toggle(name string) (string, error)
	Temperature(name string, value float64) (string, error)
}

// UsingClient is factory function to create a Fritz API interaction point.
func UsingClient(client *Client) Fritz {
	return &fritzImpl{client: client}
}

type fritzImpl struct {
	client *Client
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

// SwitchOn switches a device on. The device is identified by its name.
func (fritz *fritzImpl) SwitchOn(name string) (string, error) {
	ain, errGetAin := fritz.GetAinForName(name)
	if errGetAin != nil {
		return "", errGetAin
	}
	return fritz.switchForAin(ain, "setswitchon")
}

// SwitchOff switches a device off. The device is identified by its name.
func (fritz *fritzImpl) SwitchOff(name string) (string, error) {
	ain, errGetAin := fritz.GetAinForName(name)
	if errGetAin != nil {
		return "", errGetAin
	}
	return fritz.switchForAin(ain, "setswitchoff")
}

func (fritz *fritzImpl) switchForAin(ain, command string) (string, error) {
	resp, errSwitch := fritz.getWithAin(ain, command)
	return httpread.ReadFullyString(resp, errSwitch)
}

// GetAinForName returns the AIN corresponding to a device name.
func (fritz *fritzImpl) GetAinForName(name string) (string, error) {
	nameToAinTable, err := fritz.getNameToAinTable()
	if err != nil {
		return "", err
	}
	ain, ok := nameToAinTable[name]
	if ain == "" || !ok {
		names := stringutils.StringKeys(nameToAinTable)
		quoted := stringutils.Quote(names)
		return "", errors.New("No device found with name '" + name + "'. Available devices are " + strings.Join(quoted, ", "))
	}
	return ain, nil
}

func (fritz *fritzImpl) getNameToAinTable() (map[string]string, error) {
	devList, err := fritz.ListDevices()
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

// Toggle toggles the on/off state of a device.
func (fritz *fritzImpl) Toggle(name string) (string, error) {
	ain, errGetAin := fritz.GetAinForName(name)
	if errGetAin != nil {
		return "", errGetAin
	}
	return fritz.toggleForAin(ain)
}

// ToggleConcurrent toggles the on/off state of a device.
func (fritz *fritzImpl) ToggleConcurrent(names ...string) error {
	namesAndAins, err := fritz.getNameToAinTable()
	if err != nil {
		return err
	}
	for _, name := range names {
		if ain, ok := namesAndAins[name]; ain == "" || !ok {
			quoted := stringutils.Quote(stringutils.StringKeys(namesAndAins))
			return errors.New("No device found with name '" + name + "'. Available devices are " + strings.Join(quoted, ", "))
		}
	}

	type result struct {
		msg string
		err error
	}

	resultChannel := make(chan result, len(names))
	collectorChannel := make(chan result, len(names))
	for _, name := range names {
		ain := namesAndAins[name]
		go func(n, a string) {
			msg, err := fritz.toggleForAin(a)
			if err == nil {
				logger.Success("Successfully toggled device '"+n+"'; response was:", strings.TrimSpace(msg))
				resultChannel <- result{msg: msg, err: nil}
			} else {
				logger.Warn("Error while toggling device '" + n + "'; error was: " + err.Error())
				resultChannel <- result{msg: msg, err: fmt.Errorf("error toggling device '%s': %s", n, err.Error())}
			}
		}(name, ain)
	}

	var ops uint64
	go func() {
		for {
			res := <-resultChannel
			atomic.AddUint64(&ops, 1)
			collectorChannel <- res
			if atomic.LoadUint64(&ops) == uint64(len(names)) {
				close(resultChannel)
				close(collectorChannel)
				return
			}
		}
	}()

	errs := make([]error, 0, len(names))
	for res := range collectorChannel {
		if res.err != nil {
			errs = append(errs, res.err)
		}
	}
	if len(errs) > 0 {
		msgs := stringutils.ErrorMessages(errs)
		return errors.New("Not all devices could be toggled! Nested errors are: " + strings.Join(msgs, "; "))
	}
	return nil
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
	rounded := math.Round(doubledValue)
	response, err := fritz.getWithAinAndParam(ain, "sethkrtsoll", fmt.Sprintf("%d", rounded))
	return httpread.ReadFullyString(response, err)
}
