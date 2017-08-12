package cmd

import (
	"net"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestListInetStats  tests the command.
func TestListInetStats(t *testing.T) {
	config.Dir = "../testdata/config"
	config.Filename = "config_localhost_https_test.json"
	srv := mock.New().UnstartedServer()
	l, err := net.Listen("tcp", ":61666")
	assert.NoError(t, err)
	defer l.Close()
	srv.Listener = l
	srv.Start()
	defer srv.Close()
	err = listInetstatsCmd.RunE(listInetstatsCmd, []string{})
	assert.NoError(t, err)
}

// TestListInetStatsHasHelp ensures that the tested command provides a help text.
func TestListInetStatsHasHelp(t *testing.T) {
	assert.NotEmpty(t, listInetstatsCmd.Long)
}

// TestListInetStatsHasSynopsis ensures the tested command provides short a synopsis text.
func TestListInetStatsHasSynopsis(t *testing.T) {
	assert.NotEmpty(t, listInetstatsCmd.Short)
}

// TestFloat64ToString tests the conversion of float slice to string slice.
func TestFloat64ToString(t *testing.T) {
	fs := []float64{1.2, -12, 4.14, 9.72, 6.666666}
	transformable := float64Slice(fs)
	strs := transformable.formatFloats('f', 2)
	assert.NotNil(t, strs)
	assert.Len(t, strs, len(fs))
	assert.Equal(t, "1.20", strs[0])
	assert.Equal(t, "-12.00", strs[1])
	assert.Equal(t, "4.14", strs[2])
	assert.Equal(t, "9.72", strs[3])
	assert.Equal(t, "6.67", strs[4])
}
