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
	"io/ioutil"
	"os"
	"strings"

	"../client"
	"../utils"

	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

var numDroplets int
var noPrompt, dryRun bool

func init() {
	RootCmd.AddCommand(clusterCmd)

	clusterCmd.AddCommand(clusterCreateCmd)
	clusterCreateCmd.Flags().StringP("cloudconfig", "c", os.Getenv("HOME")+"/.cloud_config", "Directory location of cloud-config")
	clusterCreateCmd.Flags().IntVarP(&numDroplets, "droplets", "d", 3, "Number of droplets to provision for new cluster")
	//clusterCreateCmd.Flags().StringP("name", "n", "", "Name of new cluster. Required - no default")
	clusterCreateCmd.Flags().StringP("region", "r", "nyc1", "Region of new cluster")
	clusterCreateCmd.Flags().StringP("size", "s", "512mb", "RAM size of each droplet in new cluster")
	clusterCreateCmd.Flags().StringP("image", "i", "coreos-stable", "Image of each droplet in new cluster")

	clusterCmd.AddCommand(clusterDeleteCmd)
	clusterDeleteCmd.Flags().BoolVar(&noPrompt, "no-prompt", false, "Confirm culster deletion without individual droplet prompt.")
	clusterDeleteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "List cluster machines without commiting them for deletion.")
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Interface for DigitalOcean CoreOS Clusters",
	Long:  `Cluster subcommand provides a cluster-level abstraction of a DigitalOcean CoreOS deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var clusterCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new CoreOS Cluster.",
	Long: `Create a new CoreOS Cluster
	Discovery url must be provided in the config file under the key 'discovery_url'.`,
	Run: func(cmd *cobra.Command, args []string) {
		//if cmd.Flag("name").Value.String() == "" {
		//	fmt.Println("No name specified for new cluster. Aborting.")
		//	os.Exit(-1)
		//}
		if len(args) < 1 {
			fmt.Println("ERROR: Name of cluster to be create not given in command arguments. Exiting.")
			cmd.Help()
			os.Exit(-1)
		}
		clusterName := args[0]

		c, _ := client.CreateClient()
		CreateCluster(c, cmd, clusterName)
	},
}

var clusterDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a CoreOS Cluster.",
	Long: `Delete a CoreOS Cluster.
	User will be prompted for each deletion unless '-y/--yes' flag is given.
	Use '--dry-run' to list cluster machines without commiting them for deletion.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("ERROR: Name of cluster to be deleted not given in command arguments. Exiting.")
			cmd.Help()
			os.Exit(-1)
		}
		clusterName := args[0]

		c, _ := client.CreateClient()
		DeleteCluster(c, clusterName)
		//droplets := ListDroplets(c)

		//for _, drop := range droplets {
		//	if strings.HasPrefix(drop.Name, clusterName) {
		//		if dryRun {
		//			fmt.Println("Dry Run - Command would have deleted droplet:", drop.Name)
		//		} else if noPrompt {
		//			err := DeleteDroplet(c, drop.ID)
		//			if err != nil {
		//				fmt.Println(err)
		//			}
		//		} else {
		//			confirm, err := utils.AskForConfirmation("Would you like to delete droplet " + drop.Name)
		//			if err != nil {
		//				fmt.Println(err)
		//			}
		//			if confirm {
		//				err := DeleteDroplet(c, drop.ID)
		//				if err != nil {
		//					fmt.Println(err)
		//				}
		//			}
		//		}
		//	}
		//}

	},
}

func CreateCluster(client godo.Client, cmd *cobra.Command, clusterName string) {
	dropletNames := make([]string, numDroplets)
	for i := 0; i < numDroplets; i++ {
		dropletNames[i] = fmt.Sprint(clusterName, "-", i)
	}

	buf, err := ioutil.ReadFile(cmd.Flag("cloudconfig").Value.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	cloud_config := string(buf)

	droplet_keys, err := utils.ViperGetSSHKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	createRequest := &godo.DropletMultiCreateRequest{
		Names:  dropletNames,
		Region: cmd.Flag("region").Value.String(),
		Size:   cmd.Flag("size").Value.String(),
		Image: godo.DropletCreateImage{
			Slug: cmd.Flag("image").Value.String(),
		},
		SSHKeys:           droplet_keys,
		UserData:          cloud_config,
		PrivateNetworking: true,
	}

	droplet, _, err := client.Droplets.CreateMultiple(ctx, createRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(droplet)
}

func DeleteCluster(client godo.Client, clusterName string) {
	droplets := ListDroplets(client)
	for _, drop := range droplets {
		if strings.HasPrefix(drop.Name, clusterName) {
			if dryRun {
				fmt.Println("Dry Run - Command would have deleted droplet:", drop.Name)
			} else if noPrompt {
				err := DeleteDroplet(client, drop.ID)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				confirm, err := utils.AskForConfirmation("Would you like to delete droplet " + drop.Name)
				if err != nil {
					fmt.Println(err)
				}
				if confirm {
					err := DeleteDroplet(client, drop.ID)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
