package cmd

import "log"

// APIToken used to authenticate with Hetzer Cloud API
var APIToken string

func whenErrPrintAndExit(err error) {
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
