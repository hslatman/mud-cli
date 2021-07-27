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

	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/openconfig/ygot/ytypes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		filepath := args[0]
		json, err := internal.Contents(filepath)
		if err != nil {
			return errors.Wrap(err, "error reading file contents")
		}

		mud := &mudyang.Mudfile{}
		if err := mudyang.Unmarshal(json, mud); err != nil {
			return errors.Wrap(err, "can't unmarshal JSON")
		}
		// TODO: more validation options?
		options := &ytypes.LeafrefOptions{
			IgnoreMissingData: false,
			Log:               true,
		}
		if err = mud.Validate(options); err != nil {
			println(fmt.Sprintf("Error validating MUD: %v", err))
			return errors.Wrap(err, "Error validating MUD")
		}

		// TODO: some way to get more errors at once, if possible?
		// Or some nicer output.

		fmt.Println("MUD is valid")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
