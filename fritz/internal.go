package fritz

import (
	"net/http"

	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
)

// InternalAPI exposes Fritz!Box internal and undocumented API.
type InternalAPI interface {
	ListLanDevices() (*LanDevices, error)
	ListLogs() (*MessageLog, error)
	InternetStats() (*TrafficMonitoringData, error)
}

// Internal creates a Fritz/internal API from a given client.
func Internal(client *Client) InternalAPI {
	return &internalHTTP{client: client}
}

type internalHTTP struct {
	client *Client
}

// ListLogs lists the log statements produced by the FRITZ!Box.
func (internal *internalHTTP) ListLogs() (*MessageLog, error) {
	url := internal.
		query().
		query("mq_log", "logger:status/log").
		build()
	var logs MessageLog
	err := httpread.ReadFullyJSON(internal.getf(url), &logs)
	return &logs, err
}

// ListLanDevices lists the basic data of the LAN devices.
func (internal *internalHTTP) ListLanDevices() (*LanDevices, error) {
	url := internal.
		query().
		query("network", "landevice:settings/landevice/list(name,ip,mac,UID,dhcp,wlan,ethernet,active,static_dhcp,manu_name,wakeup,deleteable,source,online,speed,wlan_UIDs,auto_wakeup,guest,url,wlan_station_type,ethernet_port,wlan_show_in_monitor,plc,parental_control_abuse)").
		build()
	var devs LanDevices
	errRead := httpread.ReadFullyJSON(internal.getf(url), &devs)
	return &devs, errRead
}

// InternetStats up/downstream statistics reported by the FRITZ!Box.
func (internal *internalHTTP) InternetStats() (*TrafficMonitoringData, error) {
	url := internal.
		inetStat().
		query("useajax", "1").
		query("action", "get_graphic").
		build()
	var data []TrafficMonitoringData
	err := httpread.ReadFullyJSON(internal.getf(url), &data)
	return &data[0], err
}

func (internal *internalHTTP) query() fritzURLBuilder {
	return newURLBuilder(internal.client.Config).path(queryURI).query("sid", internal.client.SessionInfo.SID)
}

func (internal *internalHTTP) inetStat() fritzURLBuilder {
	return newURLBuilder(internal.client.Config).path(inetStatURI).query("sid", internal.client.SessionInfo.SID)
}

func (internal *internalHTTP) getf(url string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		logger.Debug("GET", url)
		return internal.client.HTTPClient.Get(url)
	}
}
