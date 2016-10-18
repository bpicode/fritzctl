package main

import (
	"log"
	"os"

	"github.com/fritzctl/fritzclient"
)

const (
	applicationName = "fritzctl"
)

func main() {
	log.Printf("Running application '%s' as %s", applicationName, os.Args[0])
	fritzClient, errCreate := fritzclient.NewClient("fritzctl.json")
	if errCreate != nil {
		log.Fatalln("Unable to create FRITZ!Box client:", errCreate)
	}
	fritzClient, err := fritzClient.Login()
	if err != nil {
		log.Fatalln("Unable to login:", err)
	}
}
