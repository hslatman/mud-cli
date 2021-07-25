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
	"crypto/x509"
	"fmt"
	"io/ioutil"
	fp "path/filepath"

	cms "github.com/github/ietf-cms"
	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/openconfig/ygot/ygot"
	"github.com/spf13/cobra"
	"go.step.sm/crypto/pemutil"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		json, err := internal.Contents(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}

		mud := &mudyang.Mudfile{}
		if err := mudyang.Unmarshal(json, mud); err != nil {
			fmt.Printf("unmarshaling JSON failed: %s\n", err)
			return
		}

		// TODO: update Mudfile with location for signature?
		if mudHasSignature(mud) {
			fmt.Printf("this MUD already has a signature available at: %s\n", *mud.Mud.MudSignature)
			return
		}

		signaturePath, err := internal.SignaturePath(filepath)
		fmt.Println("signature path: ", signaturePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonString, err := ygot.EmitJSON(mud, &ygot.EmitJSONConfig{
			Format: ygot.RFC7951,
			Indent: "  ",
			RFC7951Config: &ygot.RFC7951JSONConfig{
				AppendModuleName: true,
			},
			SkipValidation: false,
		})
		if err != nil {
			fmt.Printf("could not marshal MUD file into JSON: %s\n", err)
			return
		}

		data := []byte(jsonString)

		fmt.Println(data)

		var cert *x509.Certificate
		certFile := fp.Join(mudRootDir, "cert.pem")
		cert, err = pemutil.ReadCertificate(certFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		// TODO: provide additional information for verification,
		// amongst which are the signature file, the mud file and CA root

		// TODO: write to different location, based on signaturepath and close to the MUD file
		newSignaturePath := fp.Join(mudRootDir, "signature.der")
		der, err := ioutil.ReadFile(newSignaturePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		sd, err := cms.ParseSignedData(der)
		if err != nil {
			fmt.Println(err)
			return
		}

		pool := x509.NewCertPool()
		pool.AddCert(cert) // TODO: add different root (or system root, if trusted); optional flag/param
		options := x509.VerifyOptions{
			Roots: pool, // TODO: make this optional; now it's like a self-signed cert
		}
		if _, err := sd.VerifyDetached(data, options); err != nil {
			fmt.Println(err)
			return
		}

		// TODO: print output on how/where to store the file + signature?

		fmt.Println("success")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
