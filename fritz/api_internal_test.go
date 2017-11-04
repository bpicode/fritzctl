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

// TestInternalFritzAPI test the FRITZ API.
func TestInternalFritzAPI(t *testing.T) {
	testCases := []struct {
		dotest func(t *testing.T, internal *internal)
	}{
		{testListLanDevices},
		{testListLogs},
		{testInetStats},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(testCase.dotest).Pointer()).Name()), func(t *testing.T) {
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

			internal := NewInternal(client).(*internal)
			assert.NotNil(t, internal)
			testCase.dotest(t, internal)
		})
	}
}

func testInetStats(t *testing.T, internal *internal) {
	_, err := internal.InternetStats()
	assert.NoError(t, err)
}

func testListLanDevices(t *testing.T, internal *internal) {
	list, err := internal.ListLanDevices()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list.Network, 3)
}

func testListLogs(t *testing.T, internal *internal) {
	list, err := internal.ListLogs()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list.Messages, 7)
	for _, m := range list.Messages {
		assert.NotEmpty(t, m)
		assert.Len(t, m, 3)
	}
}
