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
	client, errCreate := NewClient("../testdata/config/config_localhost_test.yml")
	assert.NoError(t, errCreate)
	assert.NotNil(t, client)
}

// TestClientCreationNotOk ensures that an error is returned when the configuration file cannot be read.
func TestClientCreationNotOk(t *testing.T) {
	client, errCreate := NewClient("../testdata/config/ashdfashfvgashfvha.yml")
	assert.Error(t, errCreate)
	assert.Nil(t, client)
}

// TestClientLoginFailedCommunicationError tests the case (server down -> obtain challenge).
func TestClientLoginFailedCommunicationError(t *testing.T) {
	client, _ := NewClient("../testdata/config/config_localhost_test.yml")
	err := client.Login()
	assert.Error(t, err)
}

// TestClientLoginFailedSillyAnswerByServer tests the case (obtain challenge -> malformed server answer).
func TestClientLoginFailedSillyAnswerByServer(t *testing.T) {
	server, client := serverAndClient()
	defer server.Close()
	server.LoginResponse = "../testdata/config/silly.txt"
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
	client, _ := NewClient("../mock/client_config_template.yml")
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

	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false, CertificateFile: "../testdata/config/fritz.pem"}}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.NotNil(t, tlsConfig.RootCAs)

	subjects := tlsConfig.RootCAs.Subjects()
	assert.Len(t, subjects, 1)
	theOneSubj := subjects[0]
	fmt.Println("Imported x509 cert:\n", string(theOneSubj))

	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false, CertificateFile: "../testdata/config/emptyfile"}}
	cfg = config.Config{Pki: &config.Pki{SkipTLSVerify: false, CertificateFile: "../testdata/config/emptyfile"}}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)
}

// TestUtf8To16LE tests the UTF-8 to UTF-16 little endian conversion.
func TestUtf8To16LE(t *testing.T) {
	tcs := []struct {
		input, expect []byte
		name          string
	}{
		{name: "empty slice", input: []byte{}, expect: []byte{}},
		{name: "regular ascii", input: []byte("mytext"), expect: []byte{0x6d, 0x0, 0x79, 0x0, 0x74, 0x0, 0x65, 0x0, 0x78, 0x0, 0x74, 0x0}},
		{name: "emoticons and whitespace", input: []byte("😁😄\t\n"), expect: []byte{61, 216, 1, 222, 61, 216, 4, 222, 9, 0, 10, 0}},
		{name: "hello world", input: []byte("hello, world"), expect: []byte{0x68, 0x0, 0x65, 0x0, 0x6c, 0x0, 0x6c, 0x0, 0x6f, 0x0, 0x2c, 0x0, 0x20, 0x0, 0x77, 0x0, 0x6f, 0x0, 0x72, 0x0, 0x6c, 0x0, 0x64, 0x0}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, utf8To16LE(tc.input))
		})
	}
}
