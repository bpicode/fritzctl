package main

import (
	"log"
	"os"

	"github.com/bpicode/fritzctl/fatals"
	"github.com/bpicode/fritzctl/files"
	"github.com/bpicode/fritzctl/fritz"
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
