package configurer

import (
	"testing"

	"io/ioutil"

	"os"

	"github.com/stretchr/testify/assert"
)

// TestInit test the intialization phase of the interactive cli.
func TestInit(t *testing.T) {
	cli := New().(*cliConfigurer)
	cli.ApplyDefaults(Defaults())
	assert.Equal(t, cli.defaultValues, cli.userValues)
	cli.userValues.file = "/tmp/throwaway.json"
	assert.NotEqual(t, cli.defaultValues, cli.userValues)
}

// TestObtain test the user data acquisition phase of the cli.
func TestObtain(t *testing.T) {
	cli := New().(*cliConfigurer)
	cli.ApplyDefaults(Defaults())
	exfg := cli.Obtain()
	assert.NotNil(t, exfg)
}

// TestWrite test the configuration write phase of the cli.
func TestWrite(t *testing.T) {
	cli := New().(*cliConfigurer)
	extendedCfg := Defaults()
	tf, _ := ioutil.TempFile("", "test_fritzctl.json.")
	defer tf.Close()
	defer os.Remove(tf.Name())
	extendedCfg.file = tf.Name()
	cli.ApplyDefaults(extendedCfg)
	err := cli.Write()
	assert.NoError(t, err)
}

// TestWriteWithIOError test the write phase of the cli with error.
func TestWriteWithIOError(t *testing.T) {
	cli := New().(*cliConfigurer)
	extendedCfg := Defaults()
	extendedCfg.file = "/root/a/b/c/no/such/file/or/directory/cfg.json"
	cli.ApplyDefaults(extendedCfg)
	err := cli.Write()
	assert.Error(t, err)
}
