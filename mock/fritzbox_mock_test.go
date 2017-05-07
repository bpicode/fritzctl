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
	assert.True(t, r.StatusCode >= 200)
	assert.True(t, r.StatusCode < 400)
	fmt.Println(r)
	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	fmt.Println(string(body))


	r, err = client.Get(fritz.Server.URL + "/login_sid.lua?response=abdef&username=")
	assert.NoError(t, err)
	assert.True(t, r.StatusCode >= 200)
	assert.True(t, r.StatusCode < 400)
	fmt.Println(r)
	body, err = ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	fmt.Println(string(body))
}
