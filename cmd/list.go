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

	"github.com/jconard3/docore/client"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all dropets",
	Long:  `List all droplets.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := CreateClient()
		names := ListDroplets(client)
		fmt.Println(names)
	},
}

func init() {
	dropletCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func ListDroplets(client godo.Client) []string {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	var dropletNames []string
	droplets, _, _ := client.Droplets.List(opt)
	for _, element := range droplets {
		dropletNames = append(dropletNames, element.Name)
	}
	return dropletNames
}
