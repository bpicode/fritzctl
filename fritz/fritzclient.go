package fritz

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
	SID       string
}

// NewClient creates a new Client with default values.
func NewClient(configfile string) (*Client, error) {
	configPtr, err := FromFile(configfile)
	if err != nil {
		return nil, err
	}
	transportNoSslVerify := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := &http.Client{Transport: transportNoSslVerify}
	return &Client{Config: configPtr, transport: transportNoSslVerify, HTTPClient: httpClient}, nil
}

// Login tries to login into the box, obtaining the session id
func (client *Client) Login() (*Client, error) {
	sessionInfo, errObtain := client.ObtainChallenge()
	if errObtain != nil {
		return nil, errObtain
	}
	client.SessionInfo = sessionInfo
	log.Printf("FRITZ!Box challenge is %s", client.SessionInfo.Challenge)
	newSession, errSolve := client.SolveChallenge()
	if errSolve != nil {
		return nil, errSolve
	}
	client.SessionInfo = newSession
	log.Printf("FRITZ!Box challenge solved, login successful")
	return client, nil
}

func toUTF16andMD5(s string) string {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	hasher := md5.New()
	t := transform.NewWriter(hasher, enc)
	t.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// ObtainChallenge obtains the authentication challenge by the fritzbox.
func (client *Client) ObtainChallenge() (*SessionInfo, error) {
	url := client.Config.GetLoginURL()
	resp, errGet := client.HTTPClient.Get(url)
	if errGet != nil {
		return nil, errors.New("Error communicating with FRITZ!Box: " + errGet.Error())
	}
	defer resp.Body.Close()
	var sessionInfo SessionInfo
	errDecode := xml.NewDecoder(resp.Body).Decode(&sessionInfo)
	if errDecode != nil {
		return nil, errors.New("Error obtaining login challenge from FRITZ!Box: " + errDecode.Error())
	}
	return &sessionInfo, nil
}

// SolveChallenge tries to solve the authentication challenge by the fritzbox.
func (client *Client) SolveChallenge() (*SessionInfo, error) {
	challengeAndPassword := client.SessionInfo.Challenge + "-" + client.Config.Password
	challengeResponse := client.SessionInfo.Challenge + "-" + toUTF16andMD5(challengeAndPassword)
	url := client.Config.GetLoginResponseURL(challengeResponse)
	resp, errGet := client.HTTPClient.Get(url)
	if errGet != nil {
		return nil, errors.New("Error communicating with FRITZ!Box: " + errGet.Error())
	}
	defer resp.Body.Close()
	var sessionInfo SessionInfo
	errDecode := xml.NewDecoder(resp.Body).Decode(&sessionInfo)
	if errDecode != nil {
		return nil, errors.New("Error reading challenge response from FRITZ!Box: " + errDecode.Error())
	}
	if sessionInfo.SID == "0000000000000000" || sessionInfo.SID == "" {
		return nil,
			errors.New("Challenge not solved, got '" + sessionInfo.SID + "' as session id! Check login data!")
	}
	return &sessionInfo, nil
}
