package fritz

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/internal/errors"
	"github.com/bpicode/fritzctl/logger"
)

// Client encapsulates the FRITZ!Box interaction API.
type Client struct {
	Config      *config.Config // The client configuration.
	HTTPClient  *http.Client   // The HTTP client.
	SessionInfo *SessionInfo   // The current session data of the client.
}

// SessionInfo models the xml upon accessing the login endpoint.
// See also https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AVM_Technical_Note_-_Session_ID.pdf.
type SessionInfo struct {
	Challenge string `xml:"Challenge"` // A challenge provided by the FRITZ!Box.
	SID       string `xml:"SID"`       // The session id issued by the FRITZ!Box, "0000000000000000" is considered invalid/"no session".
	BlockTime string `xml:"BlockTime"` // The time that needs to expire before the next login attempt can be made.
	Rights    Rights `xml:"Rights"`    // The Rights associated withe the session.
}

// Rights wrap set of pairs (name, access-level).
type Rights struct {
	Names        []string `xml:"Name"`
	AccessLevels []string `xml:"Access"`
}

// NewClient creates a new Client with values read from a config file, given by the parameter configfile.
// Deprecated: use NewClientFromConfig.
func NewClient(configfile string) (*Client, error) {
	cfg, err := config.New(configfile)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read configuration")
	}
	return NewClientFromConfig(cfg), nil
}

// NewClientFromConfig creates a new Client with the passed configuration.
func NewClientFromConfig(cfg *config.Config) *Client {
	tlsConfig := tlsConfigFrom(cfg)
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	return &Client{Config: cfg, HTTPClient: httpClient}
}

// Login tries to login into the box and obtain the session id.
func (client *Client) Login() error {
	sessionInfo, err := client.obtainChallenge()
	if err != nil {
		return errors.Wrapf(err, "unable to obtain login challenge")
	}
	client.SessionInfo = sessionInfo
	logger.Debug("FRITZ!Box challenge is", client.SessionInfo.Challenge)
	newSession, err := client.solveChallenge()
	if err != nil {
		return errors.Wrapf(err, "unable to solve login challenge")
	}
	client.SessionInfo = newSession
	logger.Info("Login successful")
	return nil
}

func (client *Client) obtainChallenge() (*SessionInfo, error) {
	url := client.Config.GetLoginURL()
	getRemote := func() (*http.Response, error) {
		return client.HTTPClient.Get(url)
	}
	var sessionInfo SessionInfo
	err := httpread.XML(getRemote, &sessionInfo)
	return &sessionInfo, err
}

func (client *Client) solveChallenge() (*SessionInfo, error) {
	solveRemote := client.solveAttempt()
	var sessionInfo SessionInfo
	err := httpread.XML(solveRemote, &sessionInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "error solving FRITZ!Box authentication challenge")
	}
	if sessionInfo.SID == "0000000000000000" || sessionInfo.SID == "" {
		return nil, fmt.Errorf("challenge not solved, got '%s' as session id, check login data", sessionInfo.SID)
	}
	return &sessionInfo, nil
}

func (client *Client) solveAttempt() func() (*http.Response, error) {
	challengeAndPassword := client.SessionInfo.Challenge + "-" + client.Config.Login.Password
	challengeResponse := client.SessionInfo.Challenge + "-" + toUTF16andMD5(challengeAndPassword)
	url := client.Config.GetLoginResponseURL(challengeResponse)
	return func() (*http.Response, error) {
		return client.HTTPClient.Get(url)
	}
}

func toUTF16andMD5(s string) string {
	utf16 := utf8To16LE([]byte(s))
	m := md5.Sum(utf16)
	return fmt.Sprintf("%x", m)
}

func tlsConfigFrom(cfg *config.Config) *tls.Config {
	caCertPool := buildCertPool(cfg)
	return &tls.Config{InsecureSkipVerify: cfg.Pki.SkipTLSVerify, RootCAs: caCertPool}
}

func buildCertPool(cfg *config.Config) *x509.CertPool {
	if cfg.Pki.SkipTLSVerify {
		return nil
	}
	caCertPool := x509.NewCertPool()
	logger.Debug("Reading certificate file", cfg.Pki.CertificateFile)
	caCert, err := ioutil.ReadFile(cfg.Pki.CertificateFile)
	if err != nil {
		logger.Debug("Using host certificates as fallback. Reason: could not read certificate file: ", err)
		return nil
	}
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		logger.Warn("Using host certificates as fallback. Reason: certificate file ", cfg.Pki.CertificateFile, " is not a valid PEM file.")
		return nil
	}
	return caCertPool
}

func (client *Client) query() fritzURLBuilder {
	return newURLBuilder(client.Config).query("sid", client.SessionInfo.SID)
}

func (client *Client) getf(url string) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		logger.Debug("GET", url)
		return client.HTTPClient.Get(url)
	}
}
