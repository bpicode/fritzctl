package cliapp

import (
	"io"
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

// TestDeviceList is a unit test
func TestDeviceList(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml")
	defer srv.Close()
	cmd, _ := list()
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}

// TestSwitchOn is a unit test
func TestSwitchOn(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml", "testdata/answer_switch_on_test")
	defer srv.Close()
	cmd, _ := switchDevice()
	i := cmd.Run([]string{"on", "My device"})
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
