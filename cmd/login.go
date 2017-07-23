package cmd

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/fritz"
)

func clientLogin() *fritz.Client {
	configFile, err := config.FindConfigFile()
	assert.NoError(err, "unable to create FRITZ!Box client:", err)
	client, err := fritz.NewClient(configFile)
	assert.NoError(err, "unable to create FRITZ!Box client:", err)
	err = client.Login()
	assert.NoError(err, "unable to login:", err)
	return client
}
