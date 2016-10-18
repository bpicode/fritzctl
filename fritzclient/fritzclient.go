package fritzclient

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Client encapsulates the FRITZ!Box interaction API
type Client struct {
	Config      *Config
	transport   *http.Transport
	HTTPClient  *http.Client
	SessionInfo *SessionInfo
}

// SessionInfo models th xml upon accessing the login endpoint
type SessionInfo struct {
	Challenge string
}

// NewClient creates a new Client with default values.
func NewClient() *Client {
	configPtr := NewConfig()
	transportNoSslVerify := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: transportNoSslVerify}
	return &Client{Config: configPtr, transport: transportNoSslVerify, HTTPClient: httpClient}
}

// Login tries to login into the box, obtaining the session id
func (client *Client) Login() (*Client, error) {
	url := client.Config.LoginURL()
	resp, err := client.HTTPClient.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var sessionInfo SessionInfo
	err = xml.Unmarshal(body, &sessionInfo)
	if err != nil {
		return nil, err
	}
	client.SessionInfo = &sessionInfo
	return client, nil
}
