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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <droplet_name>",
	Short: "Get full details of a single droplet",
	Long:  `The get subcommand retrieves all related information for a given droplet described by the provided name.`,
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
			fmt.Println("Error creating digitalocean client. Aborting.\n")
			os.Exit(-1)
		}

		my_droplet, err := GetDroplet(my_client, droplet_name)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error retriving droplet. Aborting")
			os.Exit(-1)
		}
		fmt.Println(my_droplet)
	},
}

func init() {
	dropletCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//getCmd.Flags().StringP("name", "n", "", "name of droplet to describe")

	//getCmd.Flags().IntP("id", "i", 0, "id of droplet to describe")

}

func GetDroplet(client godo.Client, droplet_name string) (*godo.Droplet, error) {
	//Check which flag was given, if neither name nor id, error out w/ usage
	droplet_id, err := utils.NameToID(client, droplet_name)
	if err != nil {
		return nil, err
	}

	droplet, _, err := client.Droplets.Get(droplet_id)
	return droplet, err
}
