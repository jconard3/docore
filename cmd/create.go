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

	//"github.com/digitalocean/godo"
	//"github.com/jconard3/docore/client"
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
	},
}

func init() {
	dropletCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringP("region", "r", "nyc1", "Region to create droplet")
	createCmd.Flags().StringP("name", "n", "", "Name to give droplet")
	createCmd.Flags().StringP("size", "s", "512mb", "Size to give droplet")
	createCmd.Flags().StringP("image", "i", "coreos-stable", "Distribution image to give droplet")
	createCmd.Flags().Bool("backups", false, "Should automated backups be enabled for droplet?")
	createCmd.Flags().Bool("ipv6", false, "Should IPv6 be enabled for droplet?")
	createCmd.Flags().Bool("private_networking", false, "Should private networking be enabled for droplet?")
	createCmd.Flags().String("user_data", "", "String of desired User Data for droplet")
	createCmd.Flags().Bool("monitoring", false, "Should droplet install DO monitoring agent?")
	createCmd.Flags().StringArrayP("volume", "v", []string{}, "String array containing UID for each Block Storage volume to be attached to droplet")
	createCmd.Flags().StringArrayP("tags", "t", []string{}, "String array containing tag names to apply to droplet after creation")

}
