package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func whenErrPrintAndExit(err error) {
	if err != nil {
		log.Printf("Error: %s\n", err)
		log.Fatal(1)
	}
}

var workerIPCommand = &cobra.Command{Use: "worker-ip <worker name>",
	Short: "Prints IP of worker",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workerName := args[0]
		ip, err := getWorkerIP(workerName, NewHetznerClient())
		whenErrPrintAndExit(err)
		fmt.Println(ip)
	}}

func getWorkerIP(workerName string, hetznerClient *HetznerClient) (string, error) {
	server, err := hetznerClient.getServerByName(workerName)

	if err != nil {
		return "", fmt.Errorf("Error occured in HetznerClient: %s", err)
	}
	return server.PublicNet.IPv4.IP.String(), nil
}

func init() {
	rootCmd.AddCommand(workerIPCommand)
}
