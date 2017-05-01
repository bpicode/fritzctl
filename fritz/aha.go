package fritz

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/concurrent"
	"github.com/bpicode/fritzctl/fritzclient"
	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
	"github.com/bpicode/fritzctl/stringutils"
)

// HomeAutomationApi API definition, guided by
// https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type HomeAutomationApi interface {
	ListDevices() (*Devicelist, error)
	ConcurrentSwitchOn(names ...string) error
	ConcurrentSwitchOff(names ...string) error
	ConcurrentToggle(names ...string) error
	SwitchOn(ain string) (string, error)
	SwitchOff(ain string) (string, error)
	Toggle(ain string) (string, error)
	ConcurrentApplyTemperature(value float64, names ...string) error
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

// ConcurrentSwitchOn switches devices on. The devices are identified by their names.
func (aha *ahaHttp) ConcurrentSwitchOn(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.SwitchOn(ain)
		}
	}, names...)
}

// SwitchOn switches a device on. The device is identified by its AIN.
func (aha *ahaHttp) SwitchOn(ain string) (string, error) {
	return aha.switchForAin(ain, "setswitchon")
}

// ConcurrentSwitchOff switches devices off. The devices are identified by their names.
func (aha *ahaHttp) ConcurrentSwitchOff(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.SwitchOff(ain)
		}
	}, names...)
}

// SwitchOff switches a device off. The device is identified by its AIN.
func (aha *ahaHttp) SwitchOff(ain string) (string, error) {
	return aha.switchForAin(ain, "setswitchoff")
}

// ConcurrentToggle toggles the on/off state of devices.
func (aha *ahaHttp) ConcurrentToggle(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.Toggle(ain)
		}
	}, names...)
}

// ConcurrentToggle toggles the on/off state of a device. The device is identified by its AIN.
func (aha *ahaHttp) Toggle(ain string) (string, error) {
	url := aha.homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "setswitchtoggle").
		build()
	return httpread.ReadFullyString(aha.getf(url))
}

// ConcurrentApplyTemperature sets the desired temperature of "HKR" devices.
func (aha *ahaHttp) ConcurrentApplyTemperature(value float64, names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.ApplyTemperature(value, ain)
		}
	}, names...)
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

func (aha *ahaHttp) homeAutoSwitch() fritzURLBuilder {
	return newURLBuilder(aha.client.Config).path(homeAutomationURI).query("sid", aha.client.SessionInfo.SID)
}

func (aha *ahaHttp) doConcurrently(workFactory func(string) func() (string, error), names ...string) error {
	targets, err := buildBacklog(aha, names, workFactory)
	if err != nil {
		return err
	}
	results := concurrent.ScatterGather(targets, genericSuccessHandler, genericErrorHandler)
	return genericResult(results)
}

func genericSuccessHandler(key, message string) concurrent.Result {
	logger.Success("Successfully processed device '" + key + "'; response was: " + strings.TrimSpace(message))
	return concurrent.Result{Msg: message, Err: nil}
}

func genericErrorHandler(key, message string, err error) concurrent.Result {
	logger.Warn("Error while processing device '" + key + "'; error was: " + err.Error())
	return concurrent.Result{Msg: message, Err: fmt.Errorf("error toggling device '%s': %s", key, err.Error())}
}

func genericResult(results []concurrent.Result) error {
	if err := truncateToOne(results); err != nil {
		return errors.New("Not all devices could be processed! Nested errors are: " + err.Error())
	}
	return nil
}

func truncateToOne(results []concurrent.Result) error {
	errs := make([]error, 0, len(results))
	for _, res := range results {
		if res.Err != nil {
			errs = append(errs, res.Err)
		}
	}
	if len(errs) > 0 {
		msgs := stringutils.ErrorMessages(errs)
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}

func buildBacklog(fritz *ahaHttp, names []string, workFactory func(string) func() (string, error)) (map[string]func() (string, error), error) {
	namesAndAins, err := fritz.getNameToAinTable()
	if err != nil {
		return nil, err
	}
	targets := make(map[string]func() (string, error))
	for _, name := range names {
		ain, ok := namesAndAins[name]
		if ain == "" || !ok {
			quoted := stringutils.Quote(stringutils.StringKeys(namesAndAins))
			return nil, errors.New("No device found with name '" + name + "'. Available devices are " + strings.Join(quoted, ", "))
		}
		targets[name] = workFactory(ain)
	}
	return targets, nil
}

func (aha *ahaHttp) getNameToAinTable() (map[string]string, error) {
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
