package config

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewParserJSON tests the parser using JSON backend.
func TestNewParserJSON(t *testing.T) {
	p := NewParser(
		InDir("../testdata/config/", "config_test.json", JSON()),
	)
	c, err := p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

// TestNewParserPathShifting tests the parser when the filename suggests a directory substructure.
func TestNewParserPathShifting(t *testing.T) {
	p := NewParser(
		InDir("../testdata", "config/config_test.json", JSON()),
	)
	c, err := p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "xxxxx", c.Password)
}

// TestNewParserYAML tests the parser using YAML backend.
func TestNewParserYAML(t *testing.T) {
	p := NewParser(
		InDir("../testdata/config/", "config_test.yml", YAML()),
	)
	c, err := p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "fritz.box", c.Host)
}

// TestNewParserFileNotFound tests the parser when a file does not exist.
func TestNewParserFileNotFound(t *testing.T) {
	p := NewParser(
		InHomeDir(user.Current, "kjewhgjgsjdbgbjnjjub.json", JSON()),
	)
	_, err := p.Parse()
	assert.Error(t, err)
}

// TestNewParserHomeDirError tests the parser a user's $HOME cannot be determined.
func TestNewParserHomeDirError(t *testing.T) {
	p := NewParser(
		InHomeDir(func() (*user.User, error) {
			return nil, errors.New("cannot determine current user")
		}, "config.json", JSON()),
	)
	_, err := p.Parse()
	assert.Error(t, err)
}

// TestNewParserFileEmpty tests the parser when a config file is empty.
func TestNewParserFileEmpty(t *testing.T) {
	f, err := ioutil.TempFile("", "TestNewParserFileEmpty_config.json")
	assert.NoError(t, err)
	defer f.Close()
	defer os.Remove(f.Name())

	tmpDir, fName := path.Split(f.Name())
	p := NewParser(
		InHomeDir(func() (*user.User, error) {
			return &user.User{
				HomeDir: tmpDir,
			}, nil
		}, fName, JSON()),
	)
	_, err = p.Parse()
	assert.Error(t, err)
}

type errReader struct {
}

// Read always fails.
func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("I always fail")
}

// TestYAMLWithError test the YAML Decode when there is an IO error.
func TestYAMLWithError(t *testing.T) {
	y := YAML()
	err := y(&errReader{}, &Config{})
	assert.Error(t, err)
}
