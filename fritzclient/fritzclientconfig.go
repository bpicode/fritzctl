package fritzclient

import "fmt"

// Config stores client configuration of your FRITZ!Box
type Config struct {
	protocol string
	host     string
	loginURL string
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{protocol: "https", host: "fritz.box", loginURL: "/login_sid.lua"}
}

// LoginURL returns the URL that is queried for the login challenge
func (config *Config) LoginURL() string {
	return fmt.Sprintf("%s://%s%s", config.protocol, config.host, config.loginURL)
}
