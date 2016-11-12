package fritz

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetWithAin unit test.
func TestGetWithAin(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient).(*fritzImpl)
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
	resp, err := fritz.SwitchOn("DER device")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

// TestAPIGetSwitchDeviceOff unit test.
func TestAPIGetSwitchDeviceOff(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	resp, err := fritz.SwitchOff("DER device")
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
	_, err := fritz.SwitchOff("DER device")
	assert.Error(t, err)
}

// TestAPISwitchDeviceOffErrorUnkownDevice unit test.
func TestAPISwitchDeviceOffErrorUnkownDevice(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.SwitchOff("DER device")
	assert.Error(t, err)
}

// TestAPISwitchDeviceOnErrorUnkownDevice unit test.
func TestAPISwitchDeviceOnErrorUnkownDevice(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	_, err := fritz.SwitchOn("DER device")
	assert.Error(t, err)
}

// TestAPIGetSwitchByAinWithError unit test.
func TestAPIGetSwitchByAinWithError(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient).(*fritzImpl)
	ts.Close()
	_, err := fritz.switchForAin("123344", "off")
	assert.Error(t, err)
}

// TestAPIToggleDevice unit test.
func TestAPIToggleDevice(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient)
	resp, err := fritz.Toggle("DER device")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

// TestAPIToggleDeviceErrorServerDownAtListingStage unit test.
func TestAPIToggleDeviceErrorServerDownAtListingStage(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	ts.Close()
	fritz := UsingClient(fritzClient)
	_, err := fritz.Toggle("DER device")
	assert.Error(t, err)
}

// TestAPIToggleDeviceErrorServerDownAtToggleStage unit test.
func TestAPIToggleDeviceErrorServerDownAtToggleStage(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient).(*fritzImpl)
	ts.Close()
	_, err := fritz.toggleForAin("DER device")
	assert.Error(t, err)
}

// TestAPISetHkr unit test.
func TestAPISetHkr(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient).(*fritzImpl)
	_, err := fritz.Temperature("DER device", 12.5)
	assert.NoError(t, err)
}

// TestAPISetHkrDevNotFound unit test.
func TestAPISetHkrDevNotFound(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient).(*fritzImpl)
	_, err := fritz.Temperature("DOES-NOT-EXIST", 12.5)
	assert.Error(t, err)
}

// TestAPISetHkrErrorServerDownAtCommandStage unit test.
func TestAPISetHkrErrorServerDownAtCommandStage(t *testing.T) {
	ts, fritzClient := serverAndClient("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer ts.Close()
	fritzClient.Login()
	fritz := UsingClient(fritzClient).(*fritzImpl)
	ts.Close()
	_, err := fritz.temperatureForAin("12345", 12.5)
	assert.Error(t, err)
}
