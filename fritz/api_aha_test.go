package fritz

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestApiUsage tests how the API feels.
func TestApiUsage(t *testing.T) {
	assertions := assert.New(t)
	u, _ := url.Parse("https://localhost:12345")
	h := NewHomeAuto(
		Credentials("", "password"),
		SkipTLSVerify(),
		Certificate([]byte{}),
		URL(u),
		AuthEndpoint("/login_sid.lua"),
	)
	assertions.NotNil(h)
}

// TestNoPanic tests that the API functions should not panic.
func TestNoPanic(t *testing.T) {
	assertions := assert.New(t)
	u, _ := url.Parse("https://localhost:12345")
	h := NewHomeAuto(
		Credentials("", "password"),
		SkipTLSVerify(),
		URL(u),
	)
	assertions.NotPanics(func() {
		h.Login()
		h.List()
		h.Temp(20.0, "dev_name")
		h.On("dev_name")
		h.Off("dev_name")
		h.Toggle("dev_name")
	})

}

// TestFritzAPI test the FRITZ API.
func TestFritzAPI(t *testing.T) {
	testCases := []struct {
		test func(t *testing.T, h HomeAuto)
	}{
		{testOn},
		{testOff},
		{testOffError},
		{testOnError},
		{testToggle},
		{testTemp},
		{testTempErrorDeviceNotFound},
		{testToggleMany},
		{testToggleError},
		{testToggleErrorDeviceNotFound},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(tc.test).Pointer()).Name()), func(t *testing.T) {
			mockFritz := mock.New().Start()
			defer mockFritz.Close()
			ha := login(mockFritz, t)
			tc.test(t, ha)
		})
	}
}

func login(mock *mock.Fritz, t *testing.T) HomeAuto {
	u, err := url.Parse(mock.Server.URL)
	assert.NoError(t, err)
	client, err := NewClient("../mock/client_config_template.yml")
	assert.NoError(t, err)
	client.Config.Net.Protocol = u.Scheme
	client.Config.Net.Host = u.Host
	err = client.Login()
	assert.NoError(t, err)
	return &homeAuto{client: client, aha: newAinBased(client)}
}

func testTemp(t *testing.T, fritz HomeAuto) {
	err := fritz.Temp(12.5, "HKR_2")
	assert.NoError(t, err)
}

func testTempErrorDeviceNotFound(t *testing.T, h HomeAuto) {
	err := h.Temp(12.5, "DOES-NOT-EXIST")
	assert.Error(t, err)
}

func testOn(t *testing.T, h HomeAuto) {
	err := h.On("SWITCH_1")
	assert.NoError(t, err)
}

func testOnError(t *testing.T, h HomeAuto) {
	err := h.On("DEVICE_THAT_DOES_NOT_EXIST")
	assert.Error(t, err)
}

func testOff(t *testing.T, h HomeAuto) {
	err := h.Off("SWITCH_2")
	assert.NoError(t, err)
}

func testOffError(t *testing.T, h HomeAuto) {
	err := h.Off("DEVICE_THAT_DOES_NOT_EXIST")
	assert.Error(t, err)
}

func testToggle(t *testing.T, h HomeAuto) {
	err := h.Toggle("SWITCH_2")
	assert.NoError(t, err)
}

func testToggleMany(t *testing.T, h HomeAuto) {
	err := h.Toggle("SWITCH_1", "SWITCH_2", "SWITCH_3")
	assert.NoError(t, err)
}

func testToggleError(t *testing.T, h HomeAuto) {
	err := h.Toggle("SWITCH_1", "SWITCH_2", "SWITCH_3", "SWITCH_4_FAILING")
	assert.Error(t, err)
}

func testToggleErrorDeviceNotFound(t *testing.T, fritz HomeAuto) {
	err := fritz.Toggle("SWITCH_1", "UNKNOWN", "SWITCH_3")
	assert.Error(t, err)
}

// TestWithServerShutDown test the FRITZ API error handling when the backend is unreachable spontaneously.
func TestWithServerShutDown(t *testing.T) {
	testCases := []struct {
		test func(t *testing.T, h HomeAuto, s *httptest.Server)
	}{
		{testOffErrorServerDown},
		{testToggleServerDown},
		{testTempServerDown},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(tc.test).Pointer()).Name()), func(t *testing.T) {
			mockFritz := mock.New().Start()
			defer mockFritz.Close()
			ha := login(mockFritz, t)
			tc.test(t, ha, mockFritz.Server)
		})
	}
}

func testOffErrorServerDown(t *testing.T, h HomeAuto, s *httptest.Server) {
	s.Close()
	err := h.Off("SWITCH_1")
	assert.Error(t, err)
}

func testToggleServerDown(t *testing.T, h HomeAuto, s *httptest.Server) {
	s.Close()
	err := h.Toggle("SWITCH_1")
	assert.Error(t, err)
}

func testTempServerDown(t *testing.T, h HomeAuto, s *httptest.Server) {
	s.Close()
	err := h.Temp(12.5, "HKR_1")
	assert.Error(t, err)
}

// TestDeviceListCaching tests behavior when caching is turned on.
func TestDeviceListCaching(t *testing.T) {
	mockFritz := mock.New().Start()
	defer mockFritz.Close()
	u, err := url.Parse(mockFritz.Server.URL)
	assert.NoError(t, err)
	h := NewHomeAuto(
		URL(u),
		Caching(true),
		SkipTLSVerify(),
	)
	err = h.Login()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			l, err := h.List()
			assert.NoError(t, err)
			assert.NotNil(t, l)
		}()
	}
	wg.Wait()
}
