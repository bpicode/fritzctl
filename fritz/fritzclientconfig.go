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
	Protocol       string `json:"protocol"`
	Host           string `json:"host"`
	LoginURL       string `json:"loginURL"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	SkipTLSVerify  bool   `json:"skipTlsVerify"`
	CerificateFile string `json:"certificateFile"`
}

// FromFile  creates a new Config by reading from a file.
func FromFile(filestr string) (*Config, error) {
	logger.Info("Reading config from", filestr)
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
