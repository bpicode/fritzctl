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
