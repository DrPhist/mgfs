package main

import (
	"log"
	"os"
)

func checkError(err error, exit bool) {
	if err != nil {
		log.Fatal(err)
		if exit {
			os.Exit(1)
		}
	}
}
