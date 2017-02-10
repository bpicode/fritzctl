package config

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

// TestInit test the intialization phase of the interactive cli.
func TestInit(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	cli.ApplyDefaults(Defaults())
	assert.Equal(t, cli.defaultValues, cli.userValues)
	cli.userValues.file = "/tmp/throwaway.json"
	assert.NotEqual(t, cli.defaultValues, cli.userValues)
}

// TestObtain test the user data acquisition phase of the cli.
func TestObtain(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	cli.ApplyDefaults(Defaults())
	exfg := cli.Obtain()
	assert.NotNil(t, exfg)
}

// TestWrite test the configuration write phase of the cli.
func TestWrite(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	extendedCfg := Defaults()
	tf, _ := ioutil.TempFile("", "test_fritzctl.json.")
	defer tf.Close()
	defer os.Remove(tf.Name())
	extendedCfg.file = tf.Name()
	cli.ApplyDefaults(extendedCfg)
	err := cli.Write()
	assert.NoError(t, err)
}

// TestWriteAndRead test the configuration write with subsequent re-read.
func TestWriteAndRead(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	extendedCfg := Defaults()
	tf, _ := ioutil.TempFile("", "test_fritzctl.json.")
	defer tf.Close()
	defer os.Remove(tf.Name())
	extendedCfg.file = tf.Name()
	cli.ApplyDefaults(extendedCfg)
	err := cli.Write()
	assert.NoError(t, err)
	re, err := New(tf.Name())
	assert.NoError(t, err)
	assert.NotNil(t, re)
	assert.Equal(t, *cli.userValues.fritzCfg.Net, *re.Net)
	assert.Equal(t, *cli.userValues.fritzCfg.Login, *re.Login)
	assert.Equal(t, *cli.userValues.fritzCfg.Pki, *re.Pki)
}

// TestWriteWithIOError test the write phase of the cli with error.
func TestWriteWithIOError(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	extendedCfg := Defaults()
	extendedCfg.file = "/root/a/b/c/no/such/file/or/directory/cfg.json"
	cli.ApplyDefaults(extendedCfg)
	err := cli.Write()
	assert.Error(t, err)
}

// TestGreet tests the greeting.
func TestGreet(t *testing.T) {
	cli := NewConfigurer().(*cliConfigurer)
	assert.NotPanics(t, func() {
		cli.Greet()
	})
}
