package fritz

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"testing"

	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentFritzAPI test the FRITZ API.
func TestConcurrentFritzAPI(t *testing.T) {

	testCases := []struct {
		test func(t *testing.T, fritz *concurrentAhaHTTP)
	}{
		{testAPISwitchDeviceOn},
		{testAPISwitchDeviceOff},
		{testAPISwitchDeviceOffErrorUnknownDevice},
		{testAPISwitchDeviceOnErrorUnknownDevice},
		{testAPIToggleDevice},
		{testAPISetHkr},
		{testAPISetHkrDevNotFound},
		{testToggleConcurrent},
		{testToggleConcurrentWithOneError},
		{testToggleConcurrentWithDeviceNotFound},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(tc.test).Pointer()).Name()), func(t *testing.T) {
			mockFritz := mock.New().Start()
			defer mockFritz.Close()
			t, ha := createHaClient(mockFritz, t)
			tc.test(t, ha)
		})
	}
}

func createHaClient(mock *mock.Fritz, t *testing.T) (*testing.T, *concurrentAhaHTTP) {
	u, err := url.Parse(mock.Server.URL)
	assert.NoError(t, err)
	client, err := NewClient("../mock/client_config_template.json")
	assert.NoError(t, err)
	client.Config.Net.Protocol = u.Scheme
	client.Config.Net.Host = u.Host
	err = client.Login()
	assert.NoError(t, err)
	ha := concurrentConfigurator(HomeAutomation(client)).(*concurrentAhaHTTP)
	assert.NotNil(t, ha)
	return t, ha
}

func testAPISetHkr(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.temp(12.5, "HKR_2")
	assert.NoError(t, err)
}

func testAPISetHkrDevNotFound(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.temp(12.5, "DOES-NOT-EXIST")
	assert.Error(t, err)
}

func testAPISetHkrErrorServerDownAtCommandStage(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	server.Close()
	err := fritz.temp(12.5, "HKR_1")
	assert.Error(t, err)
}

func testAPISwitchDeviceOn(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.on("SWITCH_1")
	assert.NoError(t, err)
}

func testAPISwitchDeviceOff(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.off("SWITCH_2")
	assert.NoError(t, err)
}

func testAPISwitchDeviceOffErrorUnknownDevice(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.off("DEVICE_THAT_DOES_NOT_EXIST")
	assert.Error(t, err)
}

func testAPISwitchDeviceOnErrorUnknownDevice(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.on("DEVICE_THAT_DOES_NOT_EXIST")
	assert.Error(t, err)
}

func testAPIToggleDevice(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.toggle("SWITCH_2")
	assert.NoError(t, err)
}

func testToggleConcurrent(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.toggle("SWITCH_1", "SWITCH_2", "SWITCH_3")
	assert.NoError(t, err)
}

func testToggleConcurrentWithOneError(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.toggle("SWITCH_1", "SWITCH_2", "SWITCH_3", "SWITCH_4_FAILING")
	assert.Error(t, err)
}

func testToggleConcurrentWithDeviceNotFound(t *testing.T, fritz *concurrentAhaHTTP) {
	err := fritz.toggle("SWITCH_1", "UNKNOWN", "SWITCH_3")
	assert.Error(t, err)
}

// TestConcurrentFritzAPIWithServerShutDown test the FRITZ API.
func TestConcurrentFritzAPIWithServerShutDown(t *testing.T) {

	testCases := []struct {
		test func(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server)
	}{
		{testAPISwitchDeviceOffErrorServerDownAtListingStage},
		{testAPIToggleDeviceErrorServerDownAtListingStage},
		{testAPISetHkrErrorServerDownAtCommandStage},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(tc.test).Pointer()).Name()), func(t *testing.T) {
			mockFritz := mock.New().Start()
			defer mockFritz.Close()
			t, ha := createHaClient(mockFritz, t)
			tc.test(t, ha, mockFritz.Server)
		})
	}
}

func testAPISwitchDeviceOffErrorServerDownAtListingStage(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	server.Close()
	err := fritz.off("SWITCH_1")
	assert.Error(t, err)
}

func testAPIToggleDeviceErrorServerDownAtListingStage(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	server.Close()
	err := fritz.toggle("SWITCH_1")
	assert.Error(t, err)
}
