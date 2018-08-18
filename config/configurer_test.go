package config

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type infiniteNewLineReader struct{}

// Read fills the buffer with newlines.
func (r infiniteNewLineReader) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = '\n'
	}
	return len(b), nil
}

// TestObtain test the user data acquisition phase of the cli.
func TestObtain(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	exfg, _ := cli.Obtain(infiniteNewLineReader{})
	assert.NotNil(t, exfg)
}

type faultyReader struct{}

// Read always fails.
func (r faultyReader) Read(b []byte) (int, error) {
	return 0, errors.New("error")
}

// TestObtainWithErrorReading test error handling.
func TestObtainWithErrorReading(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	_, err := cli.Obtain(faultyReader{})
	assert.Error(t, err)
}

// TestWrite test the configuration write phase of the cli.
func TestWrite(t *testing.T) {
	tf, err := ioutil.TempFile("", "test_fritzctl.yml.")
	assert.NoError(t, err)
	defer tf.Close()
	defer os.Remove(tf.Name())
	extendedCfg := ExtendedConfig{}
	extendedCfg.file = tf.Name()
	err = extendedCfg.Write()
	assert.NoError(t, err)
}

// TestWriteAndRead test the configuration write with subsequent re-read.
func TestWriteAndRead(t *testing.T) {
	tf, err := ioutil.TempFile("", "test_fritzctl.yml.")
	assert.NoError(t, err)
	defer tf.Close()
	defer os.Remove(tf.Name())
	extendedCfg := ExtendedConfig{fritzCfg: Config{Net: new(Net), Login: new(Login), Pki: new(Pki)}}
	extendedCfg.file = tf.Name()
	err = extendedCfg.Write()
	assert.NoError(t, err)
	re, err := New(tf.Name())
	assert.NoError(t, err)
	assert.NotNil(t, re)
	assert.Equal(t, *extendedCfg.fritzCfg.Net, *re.Net)
	assert.Equal(t, *extendedCfg.fritzCfg.Login, *re.Login)
	assert.Equal(t, *extendedCfg.fritzCfg.Pki, *re.Pki)
}

// TestWriteWithIOError test the write phase of the cli with error.
func TestWriteWithIOError(t *testing.T) {
	extendedCfg := ExtendedConfig{file: ""}
	err := extendedCfg.Write()
	assert.Error(t, err)
}

// TestGreet tests the greeting.
func TestGreet(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	assert.NotPanics(t, func() {
		cli.Greet()
	})
}

// TestDefaultConfigFileLoc asserts the default config file behavior.
func TestDefaultConfigFileLoc(t *testing.T) {
	c := cliConfigurer{}
	assert.NotEmpty(t,
		c.defaultConfigLocation(func(file string) (string, error) { return "/path/to/folder/" + file, nil }))
	assert.NotEmpty(t,
		c.defaultConfigLocation(func(file string) (string, error) { return "", errors.New("didn't work") }))
}
