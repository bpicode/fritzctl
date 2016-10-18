package fritz

import (
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

// GetSwitchList lists the switeches configured in the system.
func (fritz *Fritz) GetSwitchList() (string, error) {
	response, errHTTP := fritz.get("getswitchlist")
	if errHTTP != nil {
		return "", errHTTP
	}
	body, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		return "", errRead
	}
	return string(body), nil
}
