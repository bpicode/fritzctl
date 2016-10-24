package fritz

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAPIGetSwitchList unit test.
func TestAPIGetSwitchList(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.GetSwitchList()
	assert.NoError(t, err)
}

// TestAPIGetSwitchListErrorServerDown unit test.
func TestAPIGetSwitchListErrorServerDown(t *testing.T) {
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
	_, err := fritz.getWithAinAndParam("ain", "cmd", "x=y")
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

// TestAPIGetDeviceListErrorServerDown unit test.
func TestAPIGetDeviceListErrorServerDown(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	ts.Close()
	_, err := fritz.ListDevices()
	assert.Error(t, err)
}

// TestAPIGetSwitchDeviceOn unit test.
func TestAPIGetSwitchDeviceOn(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	resp, err := fritz.Switch("DER device", "on")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

// TestAPIGetSwitchDeviceOff unit test.
func TestAPIGetSwitchDeviceOff(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	resp, err := fritz.Switch("DER device", "off")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

// TestAPIGetSwitchDeviceErrorServerDownAtListingStage unit test.
func TestAPIGetSwitchDeviceErrorServerDownAtListingStage(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	ts.Close()
	fritz := UsingClient(fritzClient)
	_, err := fritz.Switch("DER device", "off")
	assert.Error(t, err)
}

// TestAPIGetSwitchDeviceErrorUnkownDevice unit test.
func TestAPIGetSwitchDeviceErrorUnkownDevice(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.Switch("DER device", "off")
	assert.Error(t, err)
}

// TestAPIGetSwitchByAinWithError unit test.
func TestAPIGetSwitchByAinWithError(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	ts.Close()
	_, err := fritz.switchForAin("123344", "off")
	assert.Error(t, err)
}
