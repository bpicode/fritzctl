package cmd

import (
	"net"
	"net/http/httptest"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/stretchr/testify/assert"
)

// Test_certExport test the certificate export.
func Test_certExport(t *testing.T) {
	oldPlaces := defaultConfigPlaces
	defer func() { defaultConfigPlaces = oldPlaces }()
	defaultConfigPlaces = append([]config.Place{config.InDir("../testdata/config", "config_localhost_https_test.json", config.JSON())}, defaultConfigPlaces...)
	assertions := assert.New(t)
	server := httptest.NewUnstartedServer(nil)
	var err error
	server.Listener, err = net.Listen("tcp", ":61666")
	assertions.NoError(err)
	server.StartTLS()
	defer server.Close()
	err = certExportCmd.RunE(certExportCmd, []string{})
	assertions.NoError(err)
}
