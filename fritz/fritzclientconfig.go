package fritz

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bpicode/fritzctl/logger"
)

// Config stores client configuration of your FRITZ!Box
type Config struct {
	Protocol       string `json:"protocol"`        // The protocol to use when communicating with the FRITZ!Box. "http" or "https".
	Host           string `json:"host"`            // Host name or ip address of the FRITZ!Box. In most home setups "fritz.box" can be used. Other possible formats: "192.168.2.200:8080".
	LoginURL       string `json:"loginURL"`        // The URL for the login negotiation.
	Username       string `json:"username"`        // Username to log in. In user-agnostic setups this can be left empty.
	Password       string `json:"password"`        // The password correponding to the Username.
	SkipTLSVerify  bool   `json:"skipTlsVerify"`   // Skip TLS verifcation when using https.
	CerificateFile string `json:"certificateFile"` // Points to a certifiacte file (in PEM format) to verify the integrity of the FRITZ!Box.
}

// FromFile creates a new Config by reading from a file.
func FromFile(filestr string) (*Config, error) {
	logger.Info("Reading config file", filestr)
	file, errOpen := os.Open(filestr)
	if errOpen != nil {
		return nil, errors.New("Cannot open configuration file '" + filestr + "'. Nested error is: " + errOpen.Error())
	}
	conf := Config{}
	errDecode := json.NewDecoder(file).Decode(&conf)
	if errDecode != nil {
		return nil, errors.New("Unable to parse configuration file '" + filestr + "'. Nested error is: " + errDecode.Error())
	}
	return &conf, nil
}

// GetLoginURL returns the URL that is queried for the login challenge
func (config *Config) GetLoginURL() string {
	return fmt.Sprintf("%s://%s%s", config.Protocol, config.Host, config.LoginURL)
}

// GetLoginResponseURL returns the URL that is queried for the login challenge
func (config *Config) GetLoginResponseURL(response string) string {
	return fmt.Sprintf("%s?response=%s&username=%s", config.GetLoginURL(), response, config.Username)
}
