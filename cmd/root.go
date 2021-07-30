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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//var cfgFile string
var mudDir = ".mud"
var mudRootDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mud",
	Short: "mud provides several utilities for working with MUD files",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	initDir()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mud/cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initDir() {

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// TODO: prepare full directory structure?

	mudRootDir = filepath.Join(home, mudDir)
	if !dirExists(mudRootDir) {
		err = os.MkdirAll(mudRootDir, 0700) // TODO: right permissions?
		cobra.CheckErr(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Search config in MUD directory with name ".mud-cli" (without extension).
	// 	viper.AddConfigPath(mudDir)
	// 	viper.SetConfigType("yaml")
	// 	viper.SetConfigName(".cli")
	// }

	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }
}

func dirExists(dir string) bool {
	s, err := os.Stat(dir)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return s.IsDir()
}
