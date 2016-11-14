package fritz

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bpicode/fritzctl/httpread"
	"github.com/bpicode/fritzctl/logger"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Client encapsulates the FRITZ!Box interaction API.
type Client struct {
	Config      *Config         // The client configuration.
	transport   *http.Transport // HTTP transport settings.
	HTTPClient  *http.Client    // The HTTP client.
	SessionInfo *SessionInfo    // The current session data of the client.
}

// SessionInfo models the xml upon accessing the login endpoint.
type SessionInfo struct {
	Challenge string // A challenge provided by the FRITZ!Box.
	SID       string // The session id issued by the FRITZ!Box, "0000000000000000" is considered invalid/"no session".
}

// NewClient creates a new Client with default values.
func NewClient(configfile string) (*Client, error) {
	configPtr, err := FromFile(configfile)
	if err != nil {
		return nil, err
	}
	tlsConfig := tlsConfigFrom(configPtr)
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	return &Client{Config: configPtr, transport: transport, HTTPClient: httpClient}, nil
}

// Login tries to login into the box and obtain the session id.
func (client *Client) Login() (*Client, error) {
	sessionInfo, errObtain := client.obtainChallenge()
	if errObtain != nil {
		return nil, fmt.Errorf("Unable to obtain login challenge: %s", errObtain.Error())
	}
	client.SessionInfo = sessionInfo
	logger.Info("FRITZ!Box challenge is", client.SessionInfo.Challenge)
	newSession, errSolve := client.solveChallenge()
	if errSolve != nil {
		return nil, fmt.Errorf("Unable to solve login challenge: %s", errSolve.Error())
	}
	client.SessionInfo = newSession
	logger.Info("FRITZ!Box challenge solved, login successful")
	return client, nil
}

func (client *Client) obtainChallenge() (*SessionInfo, error) {
	url := client.Config.GetLoginURL()
	resp, errGet := client.HTTPClient.Get(url)
	var sessionInfo SessionInfo
	errParse := httpread.ReadFullyXML(resp, errGet, &sessionInfo)
	return &sessionInfo, errParse
}

func (client *Client) solveChallenge() (*SessionInfo, error) {
	challengeAndPassword := client.SessionInfo.Challenge + "-" + client.Config.Password
	challengeResponse := client.SessionInfo.Challenge + "-" + toUTF16andMD5(challengeAndPassword)
	url := client.Config.GetLoginResponseURL(challengeResponse)
	resp, errGet := client.HTTPClient.Get(url)
	var sessionInfo SessionInfo
	errXML := httpread.ReadFullyXML(resp, errGet, &sessionInfo)
	if errXML != nil {
		return nil, fmt.Errorf("Error solving FRITZ!Box authentication challenge: %s", errXML.Error())
	}
	if sessionInfo.SID == "0000000000000000" || sessionInfo.SID == "" {
		return nil, fmt.Errorf("Challenge not solved, got '%s' as session id! Check login data!", sessionInfo.SID)
	}
	return &sessionInfo, nil
}

func toUTF16andMD5(s string) string {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	hasher := md5.New()
	t := transform.NewWriter(hasher, enc)
	t.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func tlsConfigFrom(config *Config) *tls.Config {
	caCertPool := buildCertPool(config)
	return &tls.Config{InsecureSkipVerify: config.SkipTLSVerify, RootCAs: caCertPool}
}

func buildCertPool(config *Config) *x509.CertPool {
	if config.SkipTLSVerify {
		return nil
	}
	caCertPool := x509.NewCertPool()
	logger.Info("Reading certificate file", config.CerificateFile)
	caCert, err := ioutil.ReadFile(config.CerificateFile)
	if err != nil {
		logger.Warn("Using host certificates. Reason: could not read certificate file: ", err)
		return nil
	}
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		logger.Warn("Using host certificates. Reason: cerificate file ", config.CerificateFile, " not a valid PEM file.")
		return nil
	}
	return caCertPool
}
