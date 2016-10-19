package fritz

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"

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
	usr, errUser := user.Current()
	if errUser != nil {
		return nil, errUser
	}
	configPtr, err := FromFile(usr.HomeDir + "/" + configfile)
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
	return &sessionInfo, nil
}

// SolveChallenge tries to solve the authentication challenge by the fritzbox.
func (client *Client) SolveChallenge() (*SessionInfo, error) {
	challengeAndPassword := client.SessionInfo.Challenge + "-" + client.Config.Password
	challengeResponse := client.SessionInfo.Challenge + "-" + toUTF16andMD5(challengeAndPassword)
	url := client.Config.GetLoginResponseURL(challengeResponse)
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
	if sessionInfo.SID == "0000000000000000" {
		return nil,
			errors.New("Challenge not solved, got '" + sessionInfo.SID + "' as session id! Check login data!")
	}
	return &sessionInfo, nil
}
