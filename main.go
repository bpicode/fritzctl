package main

import (
	"log"
	"os"

	"github.com/fritzctl/fritz"
)

const (
	applicationName = "fritzctl"
)

func main() {
	log.Printf("Running application '%s' as %s", applicationName, os.Args[0])
	fritzClient, errCreate := fritz.NewClient("fritzctl.json")
	if errCreate != nil {
		log.Fatalln("Unable to create FRITZ!Box client:", errCreate)
	}
	fritzClient, errLogin := fritzClient.Login()
	if errLogin != nil {
		log.Fatalln("Unable to login:", errLogin)
	}
	fritz := fritz.UsingClient(fritzClient)
	r, e := fritz.GetSwitchList()
	log.Println(r, e)
}
