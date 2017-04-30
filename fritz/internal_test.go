package fritz

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/bpicode/fritzctl/fritzclient"
	"github.com/stretchr/testify/assert"
)

// TestInternalFritzAPI test the FRITZ API.
func TestInternalFritzAPI(t *testing.T) {

	serverAnswering := func(answers ...string) *httptest.Server {
		it := int32(-1)
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch, err := os.Open(answers[int(atomic.AddInt32(&it, 1))%len(answers)])
			defer ch.Close()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
			io.Copy(w, ch)
		}))
		return server
	}

	client := func() *fritzclient.Client {
		cl, err := fritzclient.New("../testdata/config_localhost_test.json")
		if err != nil {
			panic(err)
		}
		return cl
	}

	testCases := []struct {
		client *fritzclient.Client
		server *httptest.Server
		dotest func(t *testing.T, internalHttp *internalHttp)
	}{
		{
			client: client(),
			server: serverAnswering("../testdata/examplechallenge_sid_test.xml", "../testdata/examplechallenge_sid_test.xml", "../testdata/landevices_test.json"),
			dotest: testListLanDevices,
		},
		{
			client: client(),
			server: serverAnswering("../testdata/examplechallenge_sid_test.xml", "../testdata/examplechallenge_sid_test.xml", "../testdata/logs_test.json"),
			dotest: testListLogs,
		},
		{
			client: client(),
			server: serverAnswering("../testdata/examplechallenge_sid_test.xml", "../testdata/examplechallenge_sid_test.xml", "../testdata/traffic_mon_answer.json"),
			dotest: testInetStats,
		},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test aha api %s", runtime.FuncForPC(reflect.ValueOf(testCase.dotest).Pointer()).Name()), func(t *testing.T) {
			testCase.server.Start()
			defer testCase.server.Close()
			tsurl, err := url.Parse(testCase.server.URL)
			assert.NoError(t, err)
			testCase.client.Config.Net.Protocol = tsurl.Scheme
			testCase.client.Config.Net.Host = tsurl.Host
			loggedIn, err := testCase.client.Login()
			assert.NoError(t, err)
			internal := Internal(loggedIn).(*internalHttp)
			assert.NotNil(t, internal)
			testCase.dotest(t, internal)
		})
	}
}

func testInetStats(t *testing.T, internalHttp *internalHttp) {
	_, err := internalHttp.InternetStats()
	assert.NoError(t, err)
}

func testListLanDevices(t *testing.T, internalHttp *internalHttp) {
	list, err := internalHttp.ListLanDevices()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list.Network, 3)
}

func testListLogs(t *testing.T, internalHttp *internalHttp) {
	list, err := internalHttp.ListLogs()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list.Messages, 6)
	for _, m := range list.Messages {
		assert.NotEmpty(t, m)
		assert.Len(t, m, 3)
	}
}
