package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogin tests the mocked fritz server.
func TestLogin(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	client := http.Client{}
	r, err := client.Get(fritz.Server.URL + "/login_sid.lua")
	assert.NoError(t, err)
	assert2xxResponse(t, r)

	r, err = client.Get(fritz.Server.URL + "/login_sid.lua?response=abdef&username=")
	assert.NoError(t, err)
	assert2xxResponse(t, r)
}

// TestDeviceList tests the mocked fritz server.
func TestDeviceList(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, err := (&http.Client{}).Get(fritz.Server.URL + "/webservices/homeautoswitch.lua?switchcmd=getdevicelistinfos")
	assert.NoError(t, err)
	assert2xxResponse(t, r)
}

// TestSwitchingOn tests the mocked fritz server.
func TestSwitchingOn(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/webservices/homeautoswitch.lua?switchcmd=setswitchon")
	assert2xxResponse(t, r)
}

// TestSwitchingOff tests the mocked fritz server.
func TestSwitchingOff(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/webservices/homeautoswitch.lua?switchcmd=setswitchoff")
	assert2xxResponse(t, r)
}

// TestSwitchToggle tests the mocked fritz server.
func TestSwitchToggle(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/webservices/homeautoswitch.lua?switchcmd=setswitchtoggle")
	assert2xxResponse(t, r)
}

// TestSetTemp tests the mocked fritz server.
func TestSetTemp(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/webservices/homeautoswitch.lua?switchcmd=sethkrtsoll")
	assert2xxResponse(t, r)
}

// TestSwitchToggleFail tests the mocked fritz server.
func TestSwitchToggleFail(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/webservices/homeautoswitch.lua?switchcmd=setswitchtoggle&ain=faily")
	assert5xxResponse(t, r)
}

// TestLogs tests the mocked fritz server.
func TestLogs(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/query.lua?mq_log=logger")
	assert2xxResponse(t, r)
}

// TestLanDevices tests the mocked fritz server.
func TestLanDevices(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/query.lua?network=landevice")
	assert2xxResponse(t, r)
}

// TestTraffic tests the mocked fritz server.
func TestTraffic(t *testing.T) {
	fritz := New().Start()
	defer fritz.Close()
	r, _ := (&http.Client{}).Get(fritz.Server.URL + "/internet/inetstat_monitor.lua?action=get_graphic")
	assert2xxResponse(t, r)
}

func assert2xxResponse(t *testing.T, r *http.Response) {
	assert.True(t, r.StatusCode >= 200)
	assert.True(t, r.StatusCode < 300)
	fmt.Println(r)
	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	fmt.Println(string(body))
}

func assert5xxResponse(t *testing.T, r *http.Response) {
	assert.True(t, r.StatusCode >= 500)
	fmt.Println(r)
	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	fmt.Println(string(body))
}
