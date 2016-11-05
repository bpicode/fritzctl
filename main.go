package main

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/cliapp"
)

func main() {
	c := cliapp.Create()
	_, err := c.Run()
	assert.NoError(err, "Error:", err)
}
