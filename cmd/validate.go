/*
Copyright © 2021 Herman Slatman <hslatman>

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

	"github.com/hslatman/mud-cli/internal"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates a MUD file to be formatted correctly",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		filepath := args[0]
		mudfile, err := internal.ReadMUDFileFrom(filepath)
		if err != nil {
			return fmt.Errorf("could not get contents from %q: %w", filepath, err)
		}

		err = internal.Validate(mudfile)
		if err != nil {
			return fmt.Errorf("error validating MUD file: %w", err)
		}

		// TODO: some way to get more errors at once, if possible? Or some nicer output.
		fmt.Println("MUD is valid")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
