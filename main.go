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
	fritzClient := fritzclient.NewClient()
	resp, err := fritzClient.Login()
	if err != nil {
		log.Fatalln("Unable to obtain login challenge:", err)
	}
	log.Println(resp)
}
