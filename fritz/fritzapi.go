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

// Fritz API definition, guided by
// https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type Fritz interface {
	ListDevices() (*Devicelist, error)
	ListLanDevices() (*LanDevices, error)
	ListLogs() (*MessageLog, error)
	SwitchOn(names ...string) error
	SwitchOff(names ...string) error
	Toggle(names ...string) error
	Temperature(value float64, names ...string) error
}

// New creates a Fritz API from a given client.
func New(client *fritzclient.Client) Fritz {
	return &fritzImpl{client: client}
}

type fritzImpl struct {
	client *fritzclient.Client
}

func (fritz *fritzImpl) getf(url string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		return fritz.client.HTTPClient.Get(url)
	}
}

// ListLogs lists the log statements produced by the FRITZ!Box.
func (fritz *fritzImpl) ListLogs() (*MessageLog, error) {
	url := fritz.
		query().
		query("mq_log", "logger:status/log").
		build()
	var logs MessageLog
	err := httpread.ReadFullyJSON(fritz.getf(url), &logs)
	return &logs, err
}

// ListLanDevices lists the basic data of the LAN devices.
func (fritz *fritzImpl) ListLanDevices() (*LanDevices, error) {
	url := fritz.
		query().
		query("network", "landevice:settings/landevice/list(name,ip,mac,UID,dhcp,wlan,ethernet,active,static_dhcp,manu_name,wakeup,deleteable,source,online,speed,wlan_UIDs,auto_wakeup,guest,url,wlan_station_type,ethernet_port,wlan_show_in_monitor,plc,parental_control_abuse)").
		build()
	var devs LanDevices
	errRead := httpread.ReadFullyJSON(fritz.getf(url), &devs)
	return &devs, errRead
}

// ListDevices lists the basic data of the smart home devices.
func (fritz *fritzImpl) ListDevices() (*Devicelist, error) {
	url := fritz.
		homeAutoSwitch().
		query("switchcmd", "getdevicelistinfos").
		build()
	var deviceList Devicelist
	errRead := httpread.ReadFullyXML(fritz.getf(url), &deviceList)
	return &deviceList, errRead
}

// SwitchOn switches a device on. The device is identified by its name.
func (fritz *fritzImpl) SwitchOn(names ...string) error {
	return fritz.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return fritz.switchForAin(ain, "setswitchon")
		}
	}, names...)
}

// SwitchOff switches a device off. The device is identified by its name.
func (fritz *fritzImpl) SwitchOff(names ...string) error {
	return fritz.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return fritz.switchForAin(ain, "setswitchoff")
		}
	}, names...)
}

func (fritz *fritzImpl) switchForAin(ain, command string) (string, error) {
	url := fritz.
		homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", command).
		build()
	return httpread.ReadFullyString(fritz.getf(url))
}

// Toggle toggles the on/off state of devices.
func (fritz *fritzImpl) Toggle(names ...string) error {
	return fritz.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return fritz.toggleForAin(ain)
		}
	}, names...)
}

func (fritz *fritzImpl) toggleForAin(ain string) (string, error) {
	url := fritz.
		homeAutoSwitch().
		query("ain", ain).
		query("switchcmd", "setswitchtoggle").
		build()
	return httpread.ReadFullyString(fritz.getf(url))
}

func (fritz *fritzImpl) homeAutoSwitch() fritzURLBuilder {
	return newURLBuilder(fritz.client.Config).path(homeAutomationURI).query("sid", fritz.client.SessionInfo.SID)
}

func (fritz *fritzImpl) query() fritzURLBuilder {
	return newURLBuilder(fritz.client.Config).path(queryURI).query("sid", fritz.client.SessionInfo.SID)
}

// Temperature sets the desired temperature of "HKR" devices.
func (fritz *fritzImpl) Temperature(value float64, names ...string) error {
	return fritz.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			doubledValue := 2 * value
			rounded := math.Round(doubledValue)
			url := fritz.
				homeAutoSwitch().
				query("ain", ain).
				query("switchcmd", "sethkrtsoll").
				query("param", fmt.Sprintf("%d", rounded)).
				build()
			return httpread.ReadFullyString(fritz.getf(url))
		}
	}, names...)
}

func (fritz *fritzImpl) doConcurrently(workFactory func(string) func() (string, error), names ...string) error {
	targets, err := buildBacklog(fritz, names, workFactory)
	if err != nil {
		return err
	}
	results := concurrent.ScatterGather(targets, genericSuccessHandler, genericErrorHandler)
	return genericResult(results)
}

func genericSuccessHandler(key, messsage string) concurrent.Result {
	logger.Success("Successfully processed device '" + key + "'; response was: " + strings.TrimSpace(messsage))
	return concurrent.Result{Msg: messsage, Err: nil}
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

func buildBacklog(fritz *fritzImpl, names []string, workFactory func(string) func() (string, error)) (map[string]func() (string, error), error) {
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
