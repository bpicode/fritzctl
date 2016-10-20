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

// Unit test
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

// Unit test
func TestDeviceList(t *testing.T) {
	meta.ConfigDir = "testdata"
	meta.ConfigFilename = "config_localhost_test.json"
	srv := setupServer("testdata/loginresponse_test.xml", "testdata/loginresponse_test.xml", "testdata/devicelist_test.xml")
	defer srv.Close()
	cmd, _ := list()
	i := cmd.Run([]string{})
	assert.Equal(t, 0, i)
}
