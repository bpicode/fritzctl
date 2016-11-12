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

// TestClientCreationOk unit test.
func TestClientCreationOk(t *testing.T) {
	fritzClient, errCreate := NewClient("testdata/config_test.json")
	assert.NoError(t, errCreate)
	assert.NotNil(t, fritzClient)
}

// TestClientCreationNotOk unit test.
func TestClientCreationNotOk(t *testing.T) {
	fritzClient, errCreate := NewClient("testdata/ashdfashfvgashfvha.json")
	assert.Error(t, errCreate)
	assert.Nil(t, fritzClient)
}

// TestClientLoginFailedCommunationError unit test.
func TestClientLoginFailedCommunationError(t *testing.T) {
	fritzClient, _ := NewClient("testdata/config_localhost_test.json")
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginFailedSillyAnswerByServer unit test.
func TestClientLoginFailedSillyAnswerByServer(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_silly_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeFailed unit test.
func TestClientLoginChallengeFailed(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeSuccess unit test.
func TestClientLoginChallengeSuccess(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.NoError(t, err)
}

// TestClientLoginChallengeThenDerp unit test.
func TestClientLoginChallengeThenDerp(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_silly_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

// TestClientLoginChallengeThenServerDown unit test.
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

// TestCertHandling is a unit test for the certifiacte bindings.
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
