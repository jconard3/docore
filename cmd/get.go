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

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get full details of a single droplet",
	Long:  `The get subcommand retrieves all related information for a given droplet`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		//if len(args) != 1 {
		//	fmt.Println("get called without droplet name")
		//} else {
		//	fmt.Println("get called with droplet name", args[0])
		//}
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
	getCmd.Flags().StringP("name", "n", "", "name of droplet to describe")

	getCmd.Flags().IntP("id", "i", 0, "id of droplet to describe")

}
