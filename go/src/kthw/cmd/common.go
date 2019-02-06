package cmd

import (
	"fmt"
	"os"
)

// APIToken used to authenticate with Hetzer Cloud API
var APIToken string

// WhenErrPrintAndExit when err is not nil, print the error and exit
func WhenErrPrintAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
