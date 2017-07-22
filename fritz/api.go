package fritz

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/bpicode/fritzctl/config"
)

// HomeAuto is a client for the Home Automation HTTP Interface,
// see https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type HomeAuto interface {
	Login() error
	List() (*Devicelist, error)
	On(names ...string) error
	Off(names ...string) error
	Toggle(names ...string) error
	Temp(value float64, names ...string) error
}

type homeAuto struct {
	client *Client
	aha    HomeAutomationAPI
	cAha   ConcurrentHomeAutomationAPI
}

// Login tries to authenticate against the FRITZ!Box. If not successful, an error is returned. This method should be
// called before any of the other methods unless authentication is turned off at the FRITZ!Box itself.
func (h *homeAuto) Login() error {
	_, err := h.client.Login()
	return err
}

// List fetches the devices known at the FRITZ!Box. See Devicelist for details. If the devices could not be obtained,
// an error is returned.
func (h *homeAuto) List() (*Devicelist, error) {
	return h.aha.ListDevices()
}

// On activates the given devices. Devices are identified by their name. If any of the operations does not succeed,
// an error is returned.
func (h *homeAuto) On(names ...string) error {
	return h.cAha.SwitchOn(names...)
}

// Off deactivates the given devices. Devices are identified by their name. Inverse of On.
func (h *homeAuto) Off(names ...string) error {
	return h.cAha.SwitchOff(names...)
}

// Toggle switches the state of the given devices from ON to OFF and vice versa. Devices are identified by their name.
func (h *homeAuto) Toggle(names ...string) error {
	return h.cAha.Toggle(names...)
}

// Temp applies the temperature setting to the given devices. Devices are identified by their name.
func (h *homeAuto) Temp(value float64, names ...string) error {
	return h.cAha.ApplyTemperature(value, names...)
}

// Option applies fine-grained configuration to the HomeAuto client.
type Option func(h *homeAuto)

// NewHomeAuto a HomeAuto that communicates with the FRITZ!Box by means of the Home Automation HTTP Interface.
func NewHomeAuto(options ...Option) HomeAuto {
	client := defaultClient()
	aha := HomeAutomation(client)
	cAha := ConcurrentHomeAutomation(aha)
	homeAuto := homeAuto{
		client: client,
		aha:    aha,
		cAha:   cAha,
	}
	for _, option := range options {
		option(&homeAuto)
	}
	return &homeAuto
}

// URL sets the target host of the FRITZ!Box. Note that for usual setups, the url https://fritz.box:443 works.
func URL(u *url.URL) Option {
	return func(h *homeAuto) {
		h.client.Config.Net.Host = u.Host
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
		pool.AppendCertsFromPEM(bs)
		cfg := &tls.Config{RootCAs: pool}
		transport := &http.Transport{TLSClientConfig: cfg}
		h.client.HTTPClient.Transport = transport
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
