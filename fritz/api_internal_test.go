package fritz

import (
	"fmt"
	"net/url"
	"reflect"
	"runtime"
	"testing"

	"github.com/bpicode/fritzctl/mock"
	"github.com/stretchr/testify/assert"
)

// TestInternalFritzAPI tests the FRITZ API.
func TestInternalFritzAPI(t *testing.T) {
	testCases := []struct {
		tc func(t *testing.T, internal Internal)
	}{
		{testListLanDevices},
		{testListLogs},
		{testInetStats},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(testCase.tc).Pointer()).Name()), func(t *testing.T) {
			srv := mock.New().Start()
			defer srv.Close()

			client, err := NewClient("../mock/client_config_template.json")
			assert.NoError(t, err)
			u, err := url.Parse(srv.Server.URL)
			assert.NoError(t, err)
			client.Config.Net.Protocol = u.Scheme
			client.Config.Net.Host = u.Host

			err = client.Login()
			assert.NoError(t, err)

			internal := NewInternal(client)
			assert.NotNil(t, internal)
			testCase.tc(t, internal)
		})
	}
}

func testInetStats(t *testing.T, i Internal) {
	_, err := i.InternetStats()
	assert.NoError(t, err)
}

func testListLanDevices(t *testing.T, i Internal) {
	list, err := i.ListLanDevices()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list.Network, 3)
}

func testListLogs(t *testing.T, i Internal) {
	list, err := i.ListLogs()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list.Messages, 7)
	for _, m := range list.Messages {
		assert.NotEmpty(t, m)
		assert.Len(t, m, 3)
	}
}
