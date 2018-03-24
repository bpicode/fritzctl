package main

import (
	"fmt"
	"log"
)

func assertOrFatal(val bool, f string, v ...interface{}) {
	assertTrue(log.Fatalln, val, fmt.Sprintf(f, v...))
}

func assertTrue(fat func(v ...interface{}), val bool, f string) {
	if !val {
		fat(f)
	}
}
