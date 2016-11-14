package fritz

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"

	"fmt"

	"github.com/bpicode/fritzctl/fritzclient"
	"github.com/stretchr/testify/assert"
)

// TestFritzAPI test the FRITZ API.
func TestFritzAPI(t *testing.T) {

	serverAnswering := func(answers ...string) *httptest.Server {
		it := int32(-1)
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch, err := os.Open(answers[int(atomic.AddInt32(&it, 1))%len(answers)])
			defer ch.Close()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
			io.Copy(w, ch)
		}))
		return server
	}

	client := func() *fritzclient.Client {
		cl, err := fritzclient.New("testdata/config_localhost_test.json")
		if err != nil {
			panic(err)
		}
		return cl
	}

	testCases := []struct {
		client *fritzclient.Client
		server *httptest.Server
		dotest func(t *testing.T, fritz *fritzImpl, server *httptest.Server)
	}{
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml"),
			dotest: testGetWithAin,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_sid_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml"),
			dotest: testGetDeviceList,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml"),
			dotest: testAPIGetDeviceListErrorServerDown,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPISwitchDeviceOn,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPISwitchDeviceOff,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPISwitchDeviceOffErrorServerDownAtListingStage,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml"),
			dotest: testAPISwitchDeviceOffErrorUnknownDevice,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml"),
			dotest: testAPISwitchDeviceOnErrorUnknownDevice,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_empty_test.xml"),
			dotest: testAPISwitchOffByAinWithErrorServerDown,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPIToggleDevice,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPIToggleDeviceErrorServerDownAtListingStage,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPIToggleDeviceErrorServerDownAtToggleStage,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPISetHkr,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPISetHkrDevNotFound,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testAPISetHkrErrorServerDownAtCommandStage,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test", "testdata/answer_switch_on_test", "testdata/answer_switch_on_test"),
			dotest: testToggleConcurrent,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test", "testdata/answer_switch_on_test", ""),
			dotest: testToggleConcurrentWithOneError,
		},
		{
			client: client(),
			server: serverAnswering("testdata/examplechallenge_test.xml", "testdata/examplechallenge_sid_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test"),
			dotest: testToggleConcurrentWithDeviceNotFound,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test fritz api %s", runtime.FuncForPC(reflect.ValueOf(testCase.dotest).Pointer()).Name()), func(t *testing.T) {
			testCase.server.Start()
			defer testCase.server.Close()
			tsurl, err := url.Parse(testCase.server.URL)
			assert.NoError(t, err)
			testCase.client.Config.Protocol = tsurl.Scheme
			testCase.client.Config.Host = tsurl.Host
			loggedIn, err := testCase.client.Login()
			assert.NoError(t, err)
			fritz := New(loggedIn).(*fritzImpl)
			assert.NotNil(t, fritz)
			testCase.dotest(t, fritz, testCase.server)
		})
	}
}

func testAPISetHkr(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.Temperature(12.5, "DER device")
	assert.NoError(t, err)
}

func testAPISetHkrDevNotFound(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.Temperature(12.5, "DOES-NOT-EXIST")
	assert.Error(t, err)
}

func testAPISetHkrErrorServerDownAtCommandStage(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	server.Close()
	_, err := fritz.temperatureForAin("12345", 12.5)
	assert.Error(t, err)
}

func testGetWithAin(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	_, err := fritz.getWithAinAndParam("ain", "cmd", "x=y")
	assert.NoError(t, err)
}

func testGetDeviceList(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
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

func testAPIGetDeviceListErrorServerDown(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	server.Close()
	_, err := fritz.ListDevices()
	assert.Error(t, err)
}

func testAPISwitchDeviceOn(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.SwitchOn("DER device")
	assert.NoError(t, err)
}

func testAPISwitchDeviceOff(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.SwitchOff("DER device")
	assert.NoError(t, err)
}

func testAPISwitchDeviceOffErrorServerDownAtListingStage(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	server.Close()
	err := fritz.SwitchOff("DER device")
	assert.Error(t, err)
}

func testAPISwitchDeviceOffErrorUnknownDevice(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.SwitchOff("DER device")
	assert.Error(t, err)
}

func testAPISwitchDeviceOnErrorUnknownDevice(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.SwitchOn("DER device")
	assert.Error(t, err)
}

func testAPISwitchOffByAinWithErrorServerDown(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	server.Close()
	_, err := fritz.switchForAin("123344", "off")
	assert.Error(t, err)
}

func testAPIToggleDevice(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.Toggle("DER device")
	assert.NoError(t, err)
}

func testAPIToggleDeviceErrorServerDownAtListingStage(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	server.Close()
	err := fritz.Toggle("DER device")
	assert.Error(t, err)
}

func testAPIToggleDeviceErrorServerDownAtToggleStage(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	server.Close()
	_, err := fritz.toggleForAin("DER device")
	assert.Error(t, err)
}

func testToggleConcurrent(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.Toggle("DER device", "My device", "My other device")
	assert.NoError(t, err)
}

func testToggleConcurrentWithOneError(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.Toggle("DER device", "My device", "My other device")
	assert.Error(t, err)
}

func testToggleConcurrentWithDeviceNotFound(t *testing.T, fritz *fritzImpl, server *httptest.Server) {
	err := fritz.Toggle("DER device", "UNKNOWN", "My other device")
	assert.Error(t, err)
}
