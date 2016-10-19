package fritz

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientCreationOk(t *testing.T) {
	fritzClient, errCreate := NewClient("testdata/config_test.json")
	assert.NoError(t, errCreate)
	assert.NotNil(t, fritzClient)
}

func TestClientCreationNotOk(t *testing.T) {
	fritzClient, errCreate := NewClient("testdata/ashdfashfvgashfvha.json")
	assert.Error(t, errCreate)
	assert.Nil(t, fritzClient)
}

func TestClientLoginFailedCommunationError(t *testing.T) {
	fritzClient, _ := NewClient("testdata/config_localhost_test.json")
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

func TestClientLoginFailedSillyAnswerByServer(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_silly_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

func TestClientLoginChallengeFailed(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

func TestClientLoginChallengeSuccess(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.NoError(t, err)
}

func TestClientLoginChallengeThenDerp(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_silly_test.xml")
	defer ts.Close()
	_, err := fritzClient.Login()
	assert.Error(t, err)
}

func TestClientLoginChallengeThenServerDown(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml")
	defer ts.Close()

	session, errObtain := fritzClient.ObtainChallenge()
	fritzClient.SessionInfo = session
	assert.NoError(t, errObtain)

	ts.Close()
	_, err := fritzClient.SolveChallenge()
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
