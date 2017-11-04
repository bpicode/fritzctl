package fritz

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/logger"
)

// URL sets the target host of the FRITZ!Box. Note that for usual setups, the url https://fritz.box:443 works.
func URL(u *url.URL) Option {
	return func(h *homeAuto) {
		h.client.Config.Net.Host = u.Hostname()
		h.client.Config.Net.Port = u.Port()
		h.client.Config.Net.Protocol = u.Scheme
	}
}

// Credentials configures the username and password for authentication. If one wants to use the default admin account,
// the username should be an empty string.
func Credentials(username, password string) Option {
	return func(h *homeAuto) {
		h.client.Config.Login.Username = username
		h.client.Config.Login.Password = password
	}
}

// SkipTLSVerify omits TLS verification of the FRITZ!Box server. It is not recommended to use it, rather go for the
// an explicit option with Certificate.
func SkipTLSVerify() Option {
	return func(h *homeAuto) {
		skipTLS := &tls.Config{InsecureSkipVerify: true}
		h.client.HTTPClient.Transport = &http.Transport{TLSClientConfig: skipTLS}
	}
}

// Certificate actives TLS verification of the FRITZ!Box server, where the certificate is explicitly specified as byte
// array, encoded in PEM format.
func Certificate(bs []byte) Option {
	return func(h *homeAuto) {
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(bs); !ok {
			logger.Warn("Using host certificates as fallback. Supplied certificate could not be parsed.")
		}
		cfg := &tls.Config{RootCAs: pool}
		transport := &http.Transport{TLSClientConfig: cfg}
		h.client.HTTPClient.Transport = transport
	}
}

// AuthEndpoint configures the the endpoint for authentication. The default is "/login_sid.lua".
func AuthEndpoint(s string) Option {
	return func(h *homeAuto) {
		h.client.Config.Login.LoginURL = s
	}
}

func defaultClient() *Client {
	return &Client{
		Config:      defaultConfig(),
		HTTPClient:  defaultHTTP(),
		SessionInfo: defaultSessionInfo(),
	}
}

func defaultSessionInfo() *SessionInfo {
	return &SessionInfo{}
}

func defaultHTTP() *http.Client {
	return &http.Client{}
}

func defaultConfig() *config.Config {
	return &config.Config{
		Net:   defaultTarget(),
		Pki:   defaultPki(),
		Login: defaultLogin(),
	}
}

func defaultLogin() *config.Login {
	return &config.Login{
		LoginURL: "/login_sid.lua",
	}
}
func defaultPki() *config.Pki {
	return &config.Pki{}
}

func defaultTarget() *config.Net {
	return &config.Net{
		Protocol: "https",
		Host:     "fritz.box",
		Port:     "443",
	}
}
