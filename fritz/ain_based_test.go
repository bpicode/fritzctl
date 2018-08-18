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

// TestAinBased test the FRITZ API.
func TestAinBased(t *testing.T) {
	serverFactory := func() *httptest.Server {
		return mock.New().UnstartedServer()
	}

	clientFactory := func() *Client {
		cl, err := NewClient("../mock/client_config_template.yml")
		assert.NoError(t, err)
		return cl
	}

	testCases := []struct {
		doTest func(t *testing.T, fritz *ainBasedClient, server *httptest.Server)
	}{
		{testListDevices},
		{testListDevicesErrorServerDown},
		{testSwitchForAinErrorServerDown},
		{testToggleErrorServerDown},
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
			ha := newAinBased(client).(*ainBasedClient)
			assert.NotNil(t, ha)
			testCase.doTest(t, ha, server)
		})
	}
}

func testListDevices(t *testing.T, fritz *ainBasedClient, _ *httptest.Server) {
	devList, err := fritz.listDevices()
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

func testListDevicesErrorServerDown(t *testing.T, fritz *ainBasedClient, server *httptest.Server) {
	server.Close()
	_, err := fritz.listDevices()
	assert.Error(t, err)
}

func testSwitchForAinErrorServerDown(t *testing.T, fritz *ainBasedClient, server *httptest.Server) {
	server.Close()
	_, err := fritz.switchForAin("123344", "off")
	assert.Error(t, err)
}

func testToggleErrorServerDown(t *testing.T, fritz *ainBasedClient, server *httptest.Server) {
	server.Close()
	_, err := fritz.toggle("DER device")
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

// TestUnacceptableTempValues asserts that temperatures outside the range of the spec are perceived as invalid.
func TestUnacceptableTempValues(t *testing.T) {
	assertions := assert.New(t)
	h := newAinBased(nil)

	_, err := h.applyTemperature(7.5, "1235")
	assertions.Error(err)

	_, err = h.applyTemperature(55, "1235")
	assertions.Error(err)
}
