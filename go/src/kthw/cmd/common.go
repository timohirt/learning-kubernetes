package cmd

import (
	"fmt"
	"os"
)

// APIToken used to authenticate with Hetzer Cloud API
var APIToken string

func whenErrPrintAndExit(err error) {
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
}
