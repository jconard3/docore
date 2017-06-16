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

func init() {
	RootCmd.AddCommand(dropletCmd)
	dropletCmd.AddCommand(dropletListCmd)

	dropletCmd.AddCommand(dropletCreateCmd)
	dropletCreateCmd.Flags().StringP("region", "r", "nyc1", "Region of droplet to be created")
	dropletCreateCmd.Flags().StringP("size", "s", "512mb", "RAM size of droplet to be created")
	dropletCreateCmd.Flags().StringP("image", "i", "coreos-stable", "Image of droplet to be created")

	dropletCmd.AddCommand(dropletDeleteCmd)
	dropletDeleteCmd.Flags().StringP("name", "n", "", "Name of droplet to be created. Required - no default")

	dropletCmd.AddCommand(dropletGetCmd)
}

var dropletCmd = &cobra.Command{
	Use:   "droplet",
	Short: "Interface for DigitalOcean Droplets",
	Long: `Droplet subcommand provides droplet-level interface with the
	provided DigitalOcean account.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dropletListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all droplets",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.CreateClient()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error creating DO client. Aborting")
			os.Exit(-1)
		}

		droplets := ListDroplets(c)
		for _, drop := range droplets {
			fmt.Println(drop.Name)
		}
	},
}

var dropletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a droplet",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("No name specified for creating droplet. Aborting.")
			os.Exit(-1)
		}

		c, err := client.CreateClient()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error creating DO client. Aborting")
			os.Exit(-1)
		}

		CreateDroplet(c, cmd, args[0])
	},
}

var dropletDeleteCmd = &cobra.Command{
	Use:   "delete <droplet_name>",
	Short: "Delete a droplet",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.CreateClient()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error creating DO client. Aborting")
			os.Exit(-1)
		}

		if len(args) < 1 {
			fmt.Println("No name specified for deleting droplet. Aborting.")
			os.Exit(-1)
		}

		id, err := utils.NameToID(c, args[0])
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to look up droplet ID with given name.")
			os.Exit(-1)
		}

		prompt := fmt.Sprintf("Are you sure you want to delete droplet %s, ID: %d ?", cmd.Flag("name").Value.String(), id)
		confirmed, err := utils.AskForConfirmation(prompt)
		if confirmed {
			if err := DeleteDroplet(c, id); err != nil {
				fmt.Println(err)
				fmt.Println("Failed to delete droplet.")
				os.Exit(-1)
			}
		}
	},
}

var dropletGetCmd = &cobra.Command{
	Use:   "get <droplet_name>",
	Short: "Get full details of a droplet",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("No name specified for retrieving droplet. Aborting.")
			os.Exit(-1)
		}

		c, err := client.CreateClient()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error creating DO client. Aborting")
			os.Exit(-1)
		}

		droplet, err := GetDroplet(c, args[0])
		if err != nil {
			fmt.Println(err)
			fmt.Println("error retrieving droplet. Aborting")
			os.Exit(-1)
		}
		fmt.Println(droplet)
	},
}

func ListDroplets(c godo.Client) []godo.Droplet {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	droplets, _, _ := c.Droplets.List(client.Ctx, opt)
	return droplets
}

func CreateDroplet(c godo.Client, cmd *cobra.Command, name string) {
	droplet_keys, err := utils.ViperGetSSHKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: cmd.Flag("region").Value.String(),
		Size:   cmd.Flag("size").Value.String(),
		Image: godo.DropletCreateImage{
			Slug: cmd.Flag("image").Value.String(),
		},
		SSHKeys: droplet_keys,
	}

	fmt.Println(createRequest)
	confirmed, _ := utils.AskForConfirmation("Are you sure you want to create this droplet")
	if !confirmed {
		os.Exit(-1)
	}

	droplet, _, err := c.Droplets.Create(client.Ctx, createRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Droplet ", droplet.Name, " created. Currently provisioning...")
}

func DeleteDroplet(c godo.Client, id int) error {
	_, err := c.Droplets.Delete(client.Ctx, id)
	return err
}

func GetDroplet(c godo.Client, name string) (*godo.Droplet, error) {
	id, err := utils.NameToID(c, name)
	if err != nil {
		return nil, err
	}

	droplet, _, err := c.Droplets.Get(client.Ctx, id)
	return droplet, err
}
