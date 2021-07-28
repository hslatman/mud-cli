/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Provides a graphical view of a MUD file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
		// TODO: check if file exists locally or download temp copy
		// TODO: serve the file locally; temporarily
		// TODO: open browser; show MUD Visualizer with the chosen MUD file
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
