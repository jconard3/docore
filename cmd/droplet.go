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
	dropletCmd.AddCommand(dropletlistCmd)

	dropletCmd.AddCommand(dropletcreateCmd)
	dropletcreateCmd.Flags().StringP("name", "n", "", "Name of droplet to be created. Required - no default")
	dropletcreateCmd.Flags().StringP("region", "r", "nyc1", "Region of droplet to be created")
	dropletcreateCmd.Flags().StringP("size", "s", "512mb", "RAM size of droplet to be created")
	dropletcreateCmd.Flags().StringP("image", "i", "coreos-stable", "Image of droplet to be created")

	dropletCmd.AddCommand(dropletdeleteCmd)
	dropletdeleteCmd.Flags().StringP("name", "n", "", "Name of droplet to be created. Required - no default")

	dropletCmd.AddCommand(dropletgetCmd)
	dropletgetCmd.Flags().StringP("name", "n", "", "Name of droplet to get retrieved. Required - no default")
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

var dropletlistCmd = &cobra.Command{
	Use:   "list",
	Short: "List all droplets",
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := client.CreateClient()
		droplets := ListDroplets(c)
		for _, drop := range droplets {
			fmt.Println(drop.Name)
		}
	},
}

var dropletcreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a droplet",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("name").Value.String() == "" {
			fmt.Println("No name specified for creating droplet. Aborting without creating droplet.")
			os.Exit(-1)
		}

		c, _ := client.CreateClient()
		CreateDroplet(c, cmd)
	},
}

var dropletdeleteCmd = &cobra.Command{
	Use:   "delete <droplet_name>",
	Short: "Delete a droplet",
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := client.CreateClient()

		if cmd.Flag("name").Value.String() == "" {
			fmt.Println("No name specified for deleting droplet. Aborting without creating droplet.")
			os.Exit(-1)
		}
		id, err := utils.NameToID(c, cmd.Flag("name").Value.String())
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

var dropletgetCmd = &cobra.Command{
	Use:   "get <droplet_name>",
	Short: "Get full details of a droplet",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("name").Value.String() == "" {
			fmt.Println("No name specified for retrieving droplet. Aborting.")
			os.Exit(-1)
		}

		c, err := client.CreateClient()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error creating DO client. Aborting")
			os.Exit(-1)
		}

		droplet, err := GetDroplet(c, cmd.Flag("name").Value.String())
		if err != nil {
			fmt.Println(err)
			fmt.Println("error retrieving droplet. Aborting")
			os.Exit(-1)
		}
		fmt.Println(droplet)
	},
}

func ListDroplets(client godo.Client) []godo.Droplet {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	droplets, _, _ := client.Droplets.List(ctx, opt)
	return droplets
}

func CreateDroplet(client godo.Client, cmd *cobra.Command) {
	droplet_keys, err := utils.ViperGetSSHKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   cmd.Flag("name").Value.String(),
		Region: cmd.Flag("region").Value.String(),
		Size:   cmd.Flag("size").Value.String(),
		Image: godo.DropletCreateImage{
			Slug: cmd.Flag("image").Value.String(),
		},
		SSHKeys: droplet_keys,
	}

	fmt.Println(createRequest)
	c, _ := utils.AskForConfirmation("Are you sure you want to create this droplet")
	if !c {
		os.Exit(-1)
	}

	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Droplet ", droplet.Name, " created. Currently provisioning...")
}

func DeleteDroplet(client godo.Client, id int) error {
	_, err := client.Droplets.Delete(ctx, id)
	return err
}

func GetDroplet(client godo.Client, name string) (*godo.Droplet, error) {
	id, err := utils.NameToID(client, name)
	if err != nil {
		return nil, err
	}

	droplet, _, err := client.Droplets.Get(ctx, id)
	return droplet, err
}
