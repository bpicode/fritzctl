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
		{testBoxInfo},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(testCase.tc).Pointer()).Name()), func(t *testing.T) {
			srv := mock.New().Start()
			defer srv.Close()
			internal := setUpClient(t, srv)
			assert.NotNil(t, internal)
			testCase.tc(t, internal)
		})
	}
}

func setUpClient(t *testing.T, srv *mock.Fritz) Internal {
	client, err := NewClient("../mock/client_config_template.yml")
	assert.NoError(t, err)
	u, err := url.Parse(srv.Server.URL)
	assert.NoError(t, err)
	client.Config.Net.Protocol = u.Scheme
	client.Config.Net.Host = u.Host
	err = client.Login()
	assert.NoError(t, err)
	internal := NewInternal(client)
	return internal
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

func testBoxInfo(t *testing.T, i Internal) {
	data, err := i.BoxInfo()
	assert.NoError(t, err)
	assert.NotZero(t, data)
}

// TestInternal_BoxInfo_WithError test the error path of BoxInfo.
func TestInternal_BoxInfo_WithError(t *testing.T) {
	m := mock.New()
	m.SystemStatus = "/path/to/nonexistent/file/jsafbjsabfjasb.html"
	srv := m.Start()
	defer srv.Close()
	internal := setUpClient(t, srv)
	_, err := internal.BoxInfo()
	assert.Error(t, err)
}
