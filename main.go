package main

import (
	"log"
	"os"

	"github.com/fritzctl/fatals"
	"github.com/fritzctl/files"
	"github.com/fritzctl/fritz"
)

const (
	applicationName = "fritzctl"
	configFilename  = "fritzctl.json"
)

func main() {
	log.Printf("Running application '%s' as %s", applicationName, os.Args[0])
	configFile, errConfigFile := files.InHomeDir(configFilename)
	fatals.AssertNoError(errConfigFile, "Unable to create FRITZ!Box client:", errConfigFile)
	fritzClient, errCreate := fritz.NewClient(configFile)
	fatals.AssertNoError(errCreate, "Unable to create FRITZ!Box client:", errCreate)
	fritzClient, errLogin := fritzClient.Login()
	fatals.AssertNoError(errLogin, "Unable to login:", errLogin)
	fritz := fritz.UsingClient(fritzClient)
	r, e := fritz.GetSwitchList()
	log.Println(r, e)
}
