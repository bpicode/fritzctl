package cmd

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/fritz"
)

func clientLogin() *fritz.Client {
	configFile, errConfigFile := config.FindConfigFile()
	assert.NoError(errConfigFile, "unable to create FRITZ!Box client:", errConfigFile)
	client, errCreate := fritz.NewClient(configFile)
	assert.NoError(errCreate, "unable to create FRITZ!Box client:", errCreate)
	client, errLogin := client.Login()
	assert.NoError(errLogin, "unable to login:", errLogin)
	return client
}
