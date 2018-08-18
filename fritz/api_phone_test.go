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

// TestPhone tests the FRITZ/Phone API.
func TestPhone(t *testing.T) {
	testCases := []struct {
		tc func(t *testing.T, phone Phone, srv *mock.Fritz)
	}{
		{testCalls},
		{testCallsNoRecords},
		{testCallsWithServerDown},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test phone api %s", runtime.FuncForPC(reflect.ValueOf(testCase.tc).Pointer()).Name()), func(t *testing.T) {
			srv := mock.New().Start()
			defer srv.Close()

			client, err := NewClient("../mock/client_config_template.yml")
			assert.NoError(t, err)
			u, err := url.Parse(srv.Server.URL)
			assert.NoError(t, err)
			client.Config.Net.Protocol = u.Scheme
			client.Config.Net.Host = u.Host

			err = client.Login()
			assert.NoError(t, err)

			phone := NewPhone(client)
			assert.NotNil(t, phone)
			testCase.tc(t, phone, srv)
		})
	}
}

func testCalls(t *testing.T, phone Phone, _ *mock.Fritz) {
	calls, err := phone.Calls()
	assert.NoError(t, err)
	assert.NotNil(t, calls)
	assert.Len(t, calls, 3)
}

func testCallsNoRecords(t *testing.T, phone Phone, srv *mock.Fritz) {
	srv.PhoneCalls = "../mock/calls_none.csv"
	calls, err := phone.Calls()
	assert.NoError(t, err)
	assert.Empty(t, calls)
}

func testCallsWithServerDown(t *testing.T, phone Phone, srv *mock.Fritz) {
	srv.Close()
	_, err := phone.Calls()
	assert.Error(t, err)
}
