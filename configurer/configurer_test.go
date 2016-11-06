package configurer

import (
	"testing"

	"io/ioutil"

	"os"

	"github.com/stretchr/testify/assert"
)

// TestInit is a unit test.
func TestInit(t *testing.T) {
	cli := CLI().(*interactiveCLI)
	cli.InitWithDefaultVaules(Defaults())
	assert.Equal(t, cli.defaultValues, cli.userValues)
	cli.userValues.file = "/tmp/throwaway.json"
	assert.NotEqual(t, cli.defaultValues, cli.userValues)
}

// TestObtain is a unit test.
func TestObtain(t *testing.T) {
	cli := CLI().(*interactiveCLI)
	cli.InitWithDefaultVaules(Defaults())
	cli.Obtain()
}

// TestWrite is a unit test.
func TestWrite(t *testing.T) {
	cli := CLI().(*interactiveCLI)
	extendedCfg := Defaults()
	tf, _ := ioutil.TempFile("", "test_fritzctl.json.")
	defer tf.Close()
	defer os.Remove(tf.Name())
	extendedCfg.file = tf.Name()
	cli.InitWithDefaultVaules(extendedCfg)
	err := cli.Write()
	assert.NoError(t, err)
}

// TestWriteWithIOError is a unit test.
func TestWriteWithIOError(t *testing.T) {
	cli := CLI().(*interactiveCLI)
	extendedCfg := Defaults()
	extendedCfg.file = "/root/a/b/c/no/such/file/or/directory/cfg.json"
	cli.InitWithDefaultVaules(extendedCfg)
	err := cli.Write()
	assert.Error(t, err)
}
