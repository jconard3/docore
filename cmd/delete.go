// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"github.com/jconard3/docore/client"
	"github.com/jconard3/docore/utils"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <droplet_name>",
	Short: "Delete a droplet",
	Long:  `Delete a droplet from digitalocean by providing the droplet's name. The command will prompt the user for a confirmation before performing any permanent action.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if len(args) < 1 {
			fmt.Println("No droplet name provided in command arguments. Aborting")
			os.Exit(-1)
		}
		droplet_name := args[0]

		my_client, err := client.CreateClient()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error creating digitalocean client. Aborting.")
			os.Exit(-1)
		}

		my_droplet_id, err := utils.NameToID(my_client, droplet_name)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to look up droplet ID with given droplet name.")
			os.Exit(-1)
		}

		prompt := fmt.Sprintf("Are you sure you want to delete droplet %s, ID: %d ?", droplet_name, my_droplet_id)
		c, err := utils.AskForConfirmation(prompt)
		if c {
			err := DeleteDroplet(my_client, my_droplet_id)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Failed to delete droplet.")
				os.Exit(-1)
			}
		}

	},
}

func init() {
	dropletCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func DeleteDroplet(client godo.Client, id int) error {
	_, err := client.Droplets.Delete(id)
	return err
}
