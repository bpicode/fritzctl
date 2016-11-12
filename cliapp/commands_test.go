package cliapp

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bpicode/fritzctl/meta"
	"github.com/stretchr/testify/assert"
)

// TestPing is a unit test
func TestPing(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml")
	defer srv.Close()
	cmd, _ := ping()
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}

func setupServer(answers ...string) *httptest.Server {
	it := 0
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch, _ := os.Open(answers[it%len(answers)])
		defer ch.Close()
		it++
		io.Copy(w, ch)
	}))
	l, _ := net.Listen("tcp", ":61666")
	server.Listener = l
	server.Start()
	return server
}

// TestSwitchesList is a unit test for listing smart home switvch devices.
func TestSwitchesList(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml")
	defer srv.Close()
	cmd, _ := listSwitches()
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}

// TestThermostatsList is a unit test for listing HKR devices
func TestThermostatsList(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml")
	defer srv.Close()
	cmd, _ := listThermostats()
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}

// TestSwitchOn is a unit test
func TestSwitchOn(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer srv.Close()
	cmd, _ := switchOnDevice()
	i := cmd.Run([]string{"My device"})
	assert.Equal(t, 0, i)
}

// TestSwitchOff is a unit test
func TestSwitchOff(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer srv.Close()
	cmd, _ := switchOffDevice()
	i := cmd.Run([]string{"My device"})
	assert.Equal(t, 0, i)
}

// TestToggle is a unit test
func TestToggle(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer srv.Close()
	cmd, _ := toggleDevice()
	i := cmd.Run([]string{"My device"})
	assert.Equal(t, 0, i)
}

// TestSetTemp is a unit test
func TestSetTemp(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer srv.Close()
	cmd, _ := temperature()
	i := cmd.Run([]string{"19.5", "My device"})
	assert.Equal(t, 0, i)
}

// TestConfigure is a unit test
func TestConfigure(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_fritzctl")
	defer os.Remove(tempDir)
	assert.NoError(t, err)
	meta.DefaultConfigDir = tempDir

	cmd, _ := configure()
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}
