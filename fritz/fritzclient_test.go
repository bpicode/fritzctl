package fritz

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestClientCreationOk ensures that no error is returned when the configuration file is read correctly.
func TestClientCreationOk(t *testing.T) {
	client, errCreate := NewClient("../testdata/config_localhost_test.json")
	assert.NoError(t, errCreate)
	assert.NotNil(t, client)
}

// TestClientCreationNotOk ensures that an error is returned when the configuration file cannot be read.
func TestClientCreationNotOk(t *testing.T) {
	client, errCreate := NewClient("../testdata/ashdfashfvgashfvha.json")
	assert.Error(t, errCreate)
	assert.Nil(t, client)
}

// TestClientLoginFailedCommunicationError tests the case (server down -> obtain challenge).
func TestClientLoginFailedCommunicationError(t *testing.T) {
	client, _ := NewClient("../testdata/config_localhost_test.json")
	err := client.Login()
	assert.Error(t, err)
}

// TestClientLoginFailedSillyAnswerByServer tests the case (obtain challenge -> malformed server answer).
func TestClientLoginFailedSillyAnswerByServer(t *testing.T) {
	server, client := serverAndClient()
	defer server.Close()
	server.LoginResponse = "../testdata/silly.txt"
	err := client.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeFailed simulates an incorrect login challenge solution.
func TestClientLoginChallengeFailed(t *testing.T) {
	server, client := serverAndClient()
	defer server.Close()
	server.LoginResponse = "../mock/login_challenge.xml" //Replay the login challenge to simulate failure.
	err := client.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeSuccess tests the regular login workflow.
func TestClientLoginChallengeSuccess(t *testing.T) {
	server, client := serverAndClient()
	defer server.Close()
	err := client.Login()
	assert.NoError(t, err)
}

// TestClientLoginChallengeThenServerDown tests the case (obtain challenge -> server down -> solve challenge).
func TestClientLoginChallengeThenServerDown(t *testing.T) {
	server, client := serverAndClient()
	defer server.Close()
	session, errObtain := client.obtainChallenge()
	client.SessionInfo = session
	assert.NoError(t, errObtain)
	server.Close()
	_, err := client.solveChallenge()
	assert.Error(t, err)
}

func serverAndClient() (*mock.Fritz, *Client) {
	f := mock.New().Start()
	u, _ := url.Parse(f.Server.URL)
	client, _ := NewClient("../mock/client_config_template.json")
	client.Config.Net.Protocol = u.Scheme
	client.Config.Net.Host = u.Host
	return f, client
}

// TestCertHandling tests the certificate bindings.
func TestCertHandling(t *testing.T) {
	cfg := config.Config{Pki: &config.Pki{SkipTLSVerify: true}}
	tlsConfig := tlsConfigFrom(&cfg)
	assert.True(t, tlsConfig.InsecureSkipVerify)

	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false}}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)

	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false, CertificateFile: "../testdata/fritz.pem"}}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.NotNil(t, tlsConfig.RootCAs)

	subjects := tlsConfig.RootCAs.Subjects()
	assert.Len(t, subjects, 1)
	theOneSubj := subjects[0]
	fmt.Println("Imported x509 cert:\n", string(theOneSubj))

	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false, CertificateFile: "../testdata/emptyfile"}}
	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false, CertificateFile: "../testdata/emptyfile"}}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)
}
