package common

import (
	"fmt"
	"os"
)

// WhenErrPrintAndExit when err is not nil, print the error and exit
func WhenErrPrintAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
