package cliapp

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritzclient"
	"github.com/bpicode/fritzctl/meta"
)

func clientLogin() *fritzclient.Client {
	configFile, errConfigFile := meta.ConfigFile()
	assert.NoError(errConfigFile, "unable to create FRITZ!Box client:", errConfigFile)
	client, errCreate := fritzclient.New(configFile)
	assert.NoError(errCreate, "unable to create FRITZ!Box client:", errCreate)
	client, errLogin := client.Login()
	assert.NoError(errLogin, "unable to login:", errLogin)
	return client
}
