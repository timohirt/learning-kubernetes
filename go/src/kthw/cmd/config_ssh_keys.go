package cmd

import (
	"fmt"
	"kthw/cmd/common"
	"log"

	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

var addSSHKeyCommand = &cobra.Command{
	Use:   "add-ssh-key <name> <file>",
	Short: "Adds a SSH public key to the config which will be set as authorized key of root user of created servers",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		file := args[1]
		err := addSSHPublicKey(name, file)
		if err != nil {
			log.Panicf("Error while adding SSH key: %s", err)
		}

		err = viper.WriteConfig()
		if err != nil {
			log.Panicf("Error while writing SSH key to config file: %s", err)
		}
		fmt.Printf("SSH key '%s' successfully added to config.\n", name)
	}}

func addSSHPublicKey(name string, file string) error {
	publicSSHKey, err := common.ParseSSHPublicKey(name, file)
	if err != nil {
		return fmt.Errorf("Error while parsing ssh key from file '%s' config: %s", file, err)
	}

	publicSSHKey.WriteToConfig()

	return nil
}
