package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiGetSwitchList(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.GetSwitchList()
	assert.NoError(t, err)
}

func TestApiGetSwitchListErrorServerDown(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	ts.Close()
	_, err := fritz.GetSwitchList()
	assert.Error(t, err)
}

func TestGetWithAin(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.getWithAin("ain", "cmd", "x=y")
	assert.NoError(t, err)
}
