package config

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReadFromFileParsesCorrectly asserts that a regular config file is digested correctly.
func TestReadFromFileParsesCorrectly(t *testing.T) {
	var config *Config
	var err error
	config, err = New("../testdata/config/config_test.yml")
	assert.NoError(t, err)
	assert.Equal(t, "fritz.box", config.Net.Host, "Host should be parsed correctly.")
	assert.Equal(t, "https", config.Net.Protocol, "Protocol should be parsed correctly.")
	assert.Equal(t, "xxxxx", config.Login.Password, "Password should be parsed correctly.")
	assert.Equal(t, "/login_sid.lua", config.Login.LoginURL, "Login URL should be parsed correctly.")
	assert.Equal(t, "", config.Login.Username, "Username should be parsed correctly.")
}

// TestConfigProducesValidLoginURL tests that the produced login URL is syntactically correct.
func TestConfigProducesValidLoginURL(t *testing.T) {
	var config *Config
	config, _ = New("../testdata/config/config_test.yml")
	loginURL := config.GetLoginURL()
	assert.NotNil(t, loginURL)
	assert.NotEmpty(t, loginURL)
	theLoginURL, err := url.Parse(loginURL)
	assert.NoError(t, err)
	assert.NotNil(t, theLoginURL)
}

// TestConfigProducesValidLoginResponseURL tests that the produced login response URL is syntactically correct.
func TestConfigProducesValidLoginResponseURL(t *testing.T) {
	var config *Config
	config, _ = New("../testdata/config/config_test.yml")
	loginResponseURL := config.GetLoginResponseURL("some-resposne")
	assert.NotNil(t, loginResponseURL)
	assert.NotEmpty(t, loginResponseURL)
	theLoginResponseURL, err := url.Parse(loginResponseURL)
	assert.NoError(t, err)
	assert.NotNil(t, theLoginResponseURL)
}

// TestReadFromFileThatDoesNotExist tests that an error is returned at the attempt to read a non-existing file.
func TestReadFromFileThatDoesNotExist(t *testing.T) {
	_, err := New("../testdata/config/21731274tjwhbfugg374t.yml")
	assert.Error(t, err)
}

// TestReadFromInvalidFile tests that an error is returned at the attempt to read a malformed config file.
func TestReadFromInvalidFile(t *testing.T) {
	_, err := New("../testdata/config/config_invalid_test.yml")
	assert.Error(t, err)
}
