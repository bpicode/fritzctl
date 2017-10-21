package main

import (
	"flag"

	"github.com/bpicode/fritzctl/mock"
)

func main() {
	flag.Parse()
	fritz := mock.New()
	fritz.Start()
	defer fritz.Close()
}
