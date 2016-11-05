package cliapp

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/meta"
)

func clientLogin() *fritz.Client {
	configFile, errConfigFile := meta.ConfigFile()
	assert.NoError(errConfigFile, "Unable to create FRITZ!Box client:", errConfigFile)
	fritzClient, errCreate := fritz.NewClient(configFile)
	assert.NoError(errCreate, "Unable to create FRITZ!Box client:", errCreate)
	fritzClient, errLogin := fritzClient.Login()
	assert.NoError(errLogin, "Unable to login:", errLogin)
	return fritzClient
}
