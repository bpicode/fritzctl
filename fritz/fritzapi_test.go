package fritz

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApiGetSwitchList unit test.
func TestApiGetSwitchList(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.GetSwitchList()
	assert.NoError(t, err)
}

// TestApiGetSwitchListErrorServerDown unit test.
func TestApiGetSwitchListErrorServerDown(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	ts.Close()
	_, err := fritz.GetSwitchList()
	assert.Error(t, err)
}

// TestGetWithAin unit test.
func TestGetWithAin(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.getWithAin("ain", "cmd", "x=y")
	assert.NoError(t, err)
}

// TestGetDeviceList unit test.
func TestGetDeviceList(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_sid_test.xml",
		"testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	devList, err := fritz.ListDevices()
	log.Println(*devList)
	assert.NoError(t, err)
	assert.NotNil(t, devList)
	assert.NotEmpty(t, devList.Devices)
	assert.NotEmpty(t, devList.Devices[0].ID)
	assert.NotEmpty(t, devList.Devices[0].Identifier)
	assert.NotEmpty(t, devList.Devices[0].Functionbitmask)
	assert.NotEmpty(t, devList.Devices[0].Fwversion)
	assert.NotEmpty(t, devList.Devices[0].Manufacturer)
	assert.Equal(t, devList.Devices[0].Present, 1)
	assert.NotEmpty(t, devList.Devices[0].Name)

}

// TestApiGetDeviceListErrorServerDown unit test.
func TestApiGetDeviceListErrorServerDown(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	ts.Close()
	_, err := fritz.ListDevices()
	assert.Error(t, err)
}
