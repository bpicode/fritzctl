package fritz

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

// TestClientCreationOk ensures that no error is returned when the configuration file is read correctly.
func TestClientCreationOk(t *testing.T) {
	fritzClient, errCreate := NewClient("testdata/config_test.json")
	assert.NoError(t, errCreate)
	assert.NotNil(t, fritzClient)
}

// TestClientCreationNotOk ensures that an error is returned when the configuration file cannot be read.
func TestClientCreationNotOk(t *testing.T) {
	fritzClient, errCreate := NewClient("testdata/ashdfashfvgashfvha.json")
	assert.Error(t, errCreate)
	assert.Nil(t, fritzClient)
}

// TestClientLoginFailedCommunationError tests the case (server down -> obtain challenge).
func TestClientLoginFailedCommunationError(t *testing.T) {
	fritzClient, _ := NewClient("testdata/config_localhost_test.json")
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginFailedSillyAnswerByServer tests the case (obtain challenge -> malformed server answer).
func TestClientLoginFailedSillyAnswerByServer(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_silly_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeFailed simulates an incorrect login challenge solution.
func TestClientLoginChallengeFailed(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeSuccess tests the regular login workflow.
func TestClientLoginChallengeSuccess(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.NoError(t, err)
}

// TestClientLoginChallengeThenDerp tests the case (obtain challenge -> solve challenge -> malformed server answer).
func TestClientLoginChallengeThenDerp(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_silly_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeThenServerDown tests the case (obtain challenge -> server down -> solve challenge).
func TestClientLoginChallengeThenServerDown(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml")
	defer ts.Close()

	session, errObtain := fritzClient.obtainChallenge()
	fritzClient.SessionInfo = session
	assert.NoError(t, errObtain)

	ts.Close()
	_, err := fritzClient.solveChallenge()
	assert.Error(t, err)
}

func serverAndClient(answers ...string) (*httptest.Server, *Client) {
	it := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch, _ := os.Open(answers[it%len(answers)])
		defer ch.Close()
		it++
		io.Copy(w, ch)
	}))

	tsurl, _ := url.Parse(server.URL)

	client, _ := NewClient("testdata/config_localhost_test.json")
	client.Config.Protocol = tsurl.Scheme
	client.Config.Host = tsurl.Host
	return server, client
}

// TestCertHandling tests the certificate bindings.
func TestCertHandling(t *testing.T) {
	cfg := Config{SkipTLSVerify: true}
	tlsConfig := tlsConfigFrom(&cfg)
	assert.True(t, tlsConfig.InsecureSkipVerify)

	cfg = Config{SkipTLSVerify: false}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)

	cfg = Config{SkipTLSVerify: false, CerificateFile: "testdata/fritz.pem"}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.NotNil(t, tlsConfig.RootCAs)

	subjs := tlsConfig.RootCAs.Subjects()
	assert.Len(t, subjs, 1)
	theOneSubj := subjs[0]
	fmt.Println("Imported x509 cert:\n", string(theOneSubj))

	cfg = Config{SkipTLSVerify: false, CerificateFile: "testdata/emptyfile"}
	tlsConfig = tlsConfigFrom(&cfg)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)

}
