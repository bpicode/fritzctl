package fritz

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/bpicode/fritzctl/fritzclient"
	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentFritzAPI test the FRITZ API.
func TestConcurrentFritzAPI(t *testing.T) {

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
		cl, err := fritzclient.New("../testdata/config_localhost_test.json")
		assert.NoError(t, err)
		return cl
	}

	testCases := []struct {
		client *fritzclient.Client
		server *httptest.Server
		dotest func(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server)
	}{
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISwitchDeviceOn,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISwitchDeviceOff,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISwitchDeviceOffErrorServerDownAtListingStage,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISwitchDeviceOffErrorUnknownDevice,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISwitchDeviceOnErrorUnknownDevice,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPIToggleDevice,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPIToggleDeviceErrorServerDownAtListingStage,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISetHkr,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISetHkrDevNotFound,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testAPISetHkrErrorServerDownAtCommandStage,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testToggleConcurrent,
		},
		{
			client: client(),
			server: serverAnswering("../testdata/examplechallenge_test.xml", "../testdata/examplechallenge_sid_test.xml", "../testdata/devicelist_test.xml", "../testdata/answer_switch_on_test", "../testdata/answer_switch_on_test", ""),
			dotest: testToggleConcurrentWithOneError,
		},
		{
			client: client(),
			server: mock.New().UnstartedServer(),
			dotest: testToggleConcurrentWithDeviceNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(tc.dotest).Pointer()).Name()), func(t *testing.T) {
			tc.server.Start()
			defer tc.server.Close()
			u, err := url.Parse(tc.server.URL)
			assert.NoError(t, err)
			tc.client.Config.Net.Protocol = u.Scheme
			tc.client.Config.Net.Host = u.Host
			loggedIn, err := tc.client.Login()
			assert.NoError(t, err)
			ha := ConcurrentHomeAutomation(HomeAutomation(loggedIn)).(*concurrentAhaHTTP)
			assert.NotNil(t, ha)
			tc.dotest(t, ha, tc.server)
		})
	}
}

func testAPISetHkr(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.ApplyTemperature(12.5, "HKR_2")
	assert.NoError(t, err)
}

func testAPISetHkrDevNotFound(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.ApplyTemperature(12.5, "DOES-NOT-EXIST")
	assert.Error(t, err)
}

func testAPISetHkrErrorServerDownAtCommandStage(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	server.Close()
	err := fritz.ApplyTemperature(12.5, "HKR_1")
	assert.Error(t, err)
}

func testAPISwitchDeviceOn(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.SwitchOn("SWITCH_1")
	assert.NoError(t, err)
}

func testAPISwitchDeviceOff(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.SwitchOff("SWITCH_2")
	assert.NoError(t, err)
}

func testAPISwitchDeviceOffErrorServerDownAtListingStage(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	server.Close()
	err := fritz.SwitchOff("SWITCH_1")
	assert.Error(t, err)
}

func testAPISwitchDeviceOffErrorUnknownDevice(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.SwitchOff("DEVICE_THAT_DOES_NOT_EXIST")
	assert.Error(t, err)
}

func testAPISwitchDeviceOnErrorUnknownDevice(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.SwitchOn("DEVICE_THAT_DOES_NOT_EXIST")
	assert.Error(t, err)
}

func testAPIToggleDevice(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.Toggle("SWITCH_2")
	assert.NoError(t, err)
}

func testAPIToggleDeviceErrorServerDownAtListingStage(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	server.Close()
	err := fritz.Toggle("SWITCH_1")
	assert.Error(t, err)
}

func testToggleConcurrent(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.Toggle("SWITCH_1", "SWITCH_2", "SWITCH_3")
	assert.NoError(t, err)
}

func testToggleConcurrentWithOneError(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.Toggle("DER device", "My device", "My other device")
	assert.Error(t, err)
}

func testToggleConcurrentWithDeviceNotFound(t *testing.T, fritz *concurrentAhaHTTP, server *httptest.Server) {
	err := fritz.Toggle("SWITCH_1", "UNKNOWN", "SWITCH_3")
	assert.Error(t, err)
}
