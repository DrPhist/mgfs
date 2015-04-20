package main

import (
	"log"
	"os"
)

// Log error if any
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkErrorAndExit(err error, status int) {
	if err != nil {
		log.Fatal(err)
		os.Exit(status)
	}
}
