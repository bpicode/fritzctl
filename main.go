package main

import (
	"log"
	"os"
)

const (
	applicationName = "fritzctl"
)

func main() {
	log.Printf("Running application '%s' as %s", applicationName, os.Args[0])
}
