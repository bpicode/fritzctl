package fritz

import (
	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/internal/errors"
)

// Internal exposes Fritz!Box internal and undocumented API.
type Internal interface {
	ListLanDevices() (*LanDevices, error)
	ListLogs() (*MessageLog, error)
	InternetStats() (*TrafficMonitoringData, error)
	BoxInfo() (*BoxData, error)
}

// NewInternal creates a Fritz/internal API from a given client.
func NewInternal(client *Client) Internal {
	return &internal{client: client}
}

type internal struct {
	client *Client
}

// ListLogs lists the log statements produced by the FRITZ!Box.
func (i *internal) ListLogs() (*MessageLog, error) {
	url := i.
		query().
		query("mq_log", "logger:status/log").
		build()
	var logs MessageLog
	err := httpread.JSON(i.client.getf(url), &logs)
	return &logs, err
}

// ListLanDevices lists the basic data of the LAN devices.
func (i *internal) ListLanDevices() (*LanDevices, error) {
	url := i.
		query().
		query("network", "landevice:settings/landevice/list(name,ip,mac,UID,dhcp,wlan,ethernet,active,wakeup,deleteable,source,online,speed,guest,url)").
		build()
	var devs LanDevices
	err := httpread.JSON(i.client.getf(url), &devs)
	return &devs, err
}

// InternetStats up/downstream statistics reported by the FRITZ!Box.
func (i *internal) InternetStats() (*TrafficMonitoringData, error) {
	url := i.
		inetStat().
		query("myXhr", "1").
		query("xhr", "1").
		query("useajax", "1").
		query("action", "get_graphic").
		build()
	var data []TrafficMonitoringData
	err := httpread.JSON(i.client.getf(url), &data)
	return &data[0], err
}

// BoxInfo queries metadata from the FRITZ!Box. Data is drawn from: https://fritz.box/cgi-bin/system_status.
func (i *internal) BoxInfo() (*BoxData, error) {
	url := i.systemStatus().build()
	h := struct {
		Body string `xml:"body"`
	}{}
	err := httpread.XML(i.client.getf(url), &h)
	if err != nil {
		return nil, errors.Wrapf(err, "could not obtain raw system status data")
	}
	return boxDataParser{}.parse(h.Body), nil
}

func (i *internal) query() fritzURLBuilder {
	return i.client.query().path(queryURI)
}

func (i *internal) inetStat() fritzURLBuilder {
	return i.client.query().path(inetStatURI)
}

func (i *internal) systemStatus() fritzURLBuilder {
	return i.client.query().path(systemStatusURI)
}
