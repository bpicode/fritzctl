package fritz

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromFileParsesCorrectly(t *testing.T) {
	var config *Config
	var err error
	config, err = FromFile("testdata/config_test.json")
	assert.NoError(t, err)
	assert.Equal(t, "fritz.box", config.Host, "Host should be parsed correctly.")
	assert.Equal(t, "https", config.Protocol, "Protocol should be parsed correctly.")
	assert.Equal(t, "xxxxx", config.Password, "Password should be parsed correctly.")
	assert.Equal(t, "/login_sid.lua", config.LoginURL, "Login URL should be parsed correctly.")
	assert.Equal(t, "", config.Username, "Username should be parsed correctly.")
}

func TestConfigProducesValidLoginURL(t *testing.T) {
	var config *Config
	config, _ = FromFile("testdata/config_test.json")
	loginURL := config.GetLoginURL()
	assert.NotNil(t, loginURL)
	assert.NotEmpty(t, loginURL)
	theLoginURL, err := url.Parse(loginURL)
	assert.NoError(t, err)
	assert.NotNil(t, theLoginURL)
}

func TestConfigProducesValidLoginResponseURL(t *testing.T) {
	var config *Config
	config, _ = FromFile("testdata/config_test.json")
	loginResponseURL := config.GetLoginResponseURL("some-resposne")
	assert.NotNil(t, loginResponseURL)
	assert.NotEmpty(t, loginResponseURL)
	theLoginResponseURL, err := url.Parse(loginResponseURL)
	assert.NoError(t, err)
	assert.NotNil(t, theLoginResponseURL)
}

func TestReadFromFileThatDoesNotExist(t *testing.T) {
	_, err := FromFile("testdata/21731274tjwhbfugg374t.json")
	assert.Error(t, err)
}

func TestReadFromInvalidFile(t *testing.T) {
	_, err := FromFile("testdata/config_invalid_test.json")
	assert.Error(t, err)
}
