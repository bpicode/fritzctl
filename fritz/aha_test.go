package fritz

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"testing"

	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestFritzAPI test the FRITZ API.
func TestFritzAPI(t *testing.T) {

	serverFactory := func() *httptest.Server {
		return mock.New().UnstartedServer()
	}

	clientFactory := func() *Client {
		cl, err := NewClient("../mock/client_config_template.json")
		assert.NoError(t, err)
		return cl
	}

	testCases := []struct {
		doTest func(t *testing.T, fritz *ahaHTTP, server *httptest.Server)
	}{
		{testGetDeviceList},
		{testAPIGetDeviceListErrorServerDown},
		{testAPISwitchOffByAinWithErrorServerDown},
		{testAPIToggleDeviceErrorServerDownAtToggleStage},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(testCase.doTest).Pointer()).Name()), func(t *testing.T) {
			server := serverFactory()
			server.Start()
			defer server.Close()
			client := clientFactory()
			u, err := url.Parse(server.URL)
			assert.NoError(t, err)
			client.Config.Net.Protocol = u.Scheme
			client.Config.Net.Host = u.Host
			err = client.Login()
			assert.NoError(t, err)
			ha := HomeAutomation(client).(*ahaHTTP)
			assert.NotNil(t, ha)
			testCase.doTest(t, ha, server)
		})
	}
}

func testGetDeviceList(t *testing.T, fritz *ahaHTTP, server *httptest.Server) {
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

func testAPIGetDeviceListErrorServerDown(t *testing.T, fritz *ahaHTTP, server *httptest.Server) {
	server.Close()
	_, err := fritz.ListDevices()
	assert.Error(t, err)
}

func testAPISwitchOffByAinWithErrorServerDown(t *testing.T, fritz *ahaHTTP, server *httptest.Server) {
	server.Close()
	_, err := fritz.switchForAin("123344", "off")
	assert.Error(t, err)
}

func testAPIToggleDeviceErrorServerDownAtToggleStage(t *testing.T, fritz *ahaHTTP, server *httptest.Server) {
	server.Close()
	_, err := fritz.Toggle("DER device")
	assert.Error(t, err)
}

// TestRounding tests rounding.
func TestRounding(t *testing.T) {
	tcs := []struct {
		expected int64
		number   float64
		name     string
	}{
		{expected: int64(1), number: 0.5, name: "round_point_five"},
		{expected: int64(0), number: 0.4, name: "round_point_four"},
		{expected: int64(0), number: 0.1, name: "round_point_one"},
		{expected: int64(0), number: -0.1, name: "round_minus_point_one"},
		{expected: int64(0), number: -0.499, name: "round_minus_point_four_nine_nine"},
		{expected: int64(156), number: 156, name: "round_one_hundred_fifty_six"},
		{expected: int64(3), number: 3.14, name: "round_pi"},
		{expected: int64(4), number: 3.54, name: "round_three_point_five_four"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, round(tc.number))
		})
	}
}
