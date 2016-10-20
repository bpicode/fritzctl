package main

import (
	"github.com/bpicode/fritzctl/cliapp"
	"github.com/bpicode/fritzctl/fatals"
)

func main() {
	c := cliapp.Create()
	_, err := c.Run()
	fatals.AssertNoError(err, "Error:", err)
}
