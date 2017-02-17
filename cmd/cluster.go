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

	"github.com/digitalocean/godo"
	"github.com/jconard3/docore/client"
	"github.com/jconard3/docore/utils"
	"github.com/spf13/cobra"
)

var numDroplets int

func init() {
	RootCmd.AddCommand(clusterCmd)

	clusterCmd.AddCommand(clustercreateCmd)
	clustercreateCmd.Flags().StringP("cloudconfig", "c", os.Getenv("HOME")+"/.cloud_config", "Directory location of cloud-config")
	clustercreateCmd.Flags().IntVarP(&numDroplets, "droplets", "d", 3, "Number of droplets to provision for new cluster")
	clustercreateCmd.Flags().StringP("name", "n", "", "Name of new cluster. Required - no default")
	clustercreateCmd.Flags().StringP("region", "r", "nyc1", "Region of new cluster")
	clustercreateCmd.Flags().StringP("size", "s", "512mb", "RAM size of each droplet in new cluster")
	clustercreateCmd.Flags().StringP("image", "i", "coreos-stable", "Image of each droplet in new cluster")
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Interface for DigitalOcean CoreOS Clusters",
	Long:  `Cluster subcommand provides a cluster-level abstraction of a DigitalOcean CoreOS deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var clustercreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new CoreOS Cluster.",
	Long: `Create a new CoreOS Cluster
	Discovery url must be provided in the config file under the key 'discovery_url'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("name").Value.String() == "" {
			fmt.Println("No name specified for new cluster. Aborting.")
			os.Exit(-1)
		}

		c, _ := client.CreateClient()
		CreateCluster(c, cmd)
	},
}

func CreateCluster(client godo.Client, cmd *cobra.Command) {
	dropletNames := make([]string, numDroplets)
	for i := 0; i < numDroplets; i++ {
		dropletNames[i] = fmt.Sprint(cmd.Flag("name").Value.String(), "-", i)
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

	droplet, _, err := client.Droplets.CreateMultiple(createRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(droplet)
}
