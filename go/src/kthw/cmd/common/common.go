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

// FileExists returns true if a given file exists.
func FileExists(filePath string) bool {
	if _, statErr := os.Stat(filePath); statErr == nil {
		return true
	}
	return false
}

// ArrayContains returns 'true' if value is in arr.
func ArrayContains(arr []string, value string) bool {
	for _, elem := range arr {
		if value == elem {
			return true
		}
	}
	return false
}
