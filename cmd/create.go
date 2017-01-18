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
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a droplet",
	Long: `Create a single droplet. 
	SSH public keys must be specified in an array of strings by the key "ssh_keys" in the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("ssh_keys") {
			fmt.Println("No ssh_keys specified in config file. Aborting without creating droplet.")
			os.Exit(-1)
		}

		if cmd.Flag("name").Value.String() == "" {
			fmt.Println("No name specified for creating droplet. Aborting without creating droplet.")
			os.Exit(-1)
		}

		c, _ := client.CreateClient()
		CreateDroplet(c, cmd)
	},
}

func init() {
	dropletCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("region", "r", "nyc1", "Region to create droplet")
	createCmd.Flags().StringP("name", "n", "", "Name to give droplet")
	createCmd.Flags().StringP("size", "s", "512mb", "Size to give droplet")
	createCmd.Flags().StringP("image", "i", "coreos-stable", "Distribution image to give droplet")
	//createCmd.Flags().Bool("backups", false, "Should automated backups be enabled for droplet?")
	//createCmd.Flags().Bool("ipv6", false, "Should IPv6 be enabled for droplet?")
	//createCmd.Flags().Bool("private_networking", false, "Should private networking be enabled for droplet?")
	//createCmd.Flags().String("user_data", "", "String of desired User Data for droplet")
	//createCmd.Flags().Bool("monitoring", false, "Should droplet install DO monitoring agent?")
	//createCmd.Flags().StringArrayP("volumes", "v", []string{}, "String array containing UID for each Block Storage volume to be attached to droplet")
	//createCmd.Flags().StringArrayP("tags", "t", []string{}, "String array containing tag names to apply to droplet after creation")
}

func CreateDroplet(client godo.Client, cmd *cobra.Command) {
	ssh_keys := viper.GetStringSlice("ssh_keys")
	droplet_keys := make([]godo.DropletCreateSSHKey, len(ssh_keys))
	for i, ssh_key := range ssh_keys {
		droplet_keys[i].Fingerprint = ssh_key
	}

	//volumes := cmd.Flag("volumes").Value.String()
	//droplet_volumes := make([]godo.DropletCreateVolume, len(volumes))
	//for i, volume := range volumes {
	//	droplet_volumes[i].ID = volume
	//}

	createRequest := &godo.DropletCreateRequest{
		Name:   cmd.Flag("name").Value.String(),
		Region: cmd.Flag("region").Value.String(),
		Size:   cmd.Flag("size").Value.String(),
		Image: godo.DropletCreateImage{
			Slug: cmd.Flag("image").Value.String(),
		},
		SSHKeys: droplet_keys,
		//Backups: cmd.Flag("backups").Value.String(),
		//	IPv6:              cmd.Flag("ipv6").Value.String(),
		//	PrivateNetworking: cmd.Flag("private_networking").Value.String(),
		//	Monitoring:        cmd.Flag("monitoring").Value.String(),
		//	UserData:          cmd.Flag("user_data").Value.String(),
		//	Volumes:           cmd.Flag("volumes").Value,
		//	Tags:              cmd.Flag("tags").Value,
	}

	fmt.Println(createRequest)
	c, _ := utils.AskForConfirmation("Are you sure you want to create this droplet")
	if !c {
		os.Exit(-1)
	}

	droplet, _, err := client.Droplets.Create(createRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Droplet ", droplet.Name, " created. Currently provisioning...")

}
