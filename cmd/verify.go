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
	"log"
	fp "path/filepath"

	cms "github.com/github/ietf-cms"
	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.step.sm/crypto/pemutil"
)

var signatureFlag string

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the signature for a MUD file",
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]
		data, err := internal.Read(filepath)
		if err != nil {
			errors.Wrapf(err, "error reading contents of %s", filepath)
		}
		mudfile, err := internal.Parse(data)
		if err != nil {
			return errors.Wrap(err, "could not get contents")
		}
		signature, err := getSignatureFilepath(mudfile)
		if err != nil {
			return errors.Wrap(err, "retrieving signature from MUD failed")
		}

		fmt.Println("signature path: ", signature)

		// signaturePath, err := internal.SignaturePath(filepath)
		// fmt.Println("signature path: ", signaturePath)
		// if err != nil {
		// 	return errors.Wrap(err, "retrieving signature path from MUD failed")
		// }

		// TODO: read MUD signature file location from MUD file, retrieve it and verify using that one
		// TODO: allow for providing the signature file as an argument

		// jsonString, err := ygot.EmitJSON(mudfile, &ygot.EmitJSONConfig{
		// 	Format: ygot.RFC7951,
		// 	Indent: "  ",
		// 	RFC7951Config: &ygot.RFC7951JSONConfig{
		// 		AppendModuleName: true,
		// 	},
		// 	SkipValidation: false,
		// })
		// if err != nil {
		// 	return errors.Wrap(err, "could not marshal MUD file into JSON")
		// }

		// data := []byte(jsonString)

		fmt.Println(data)

		var cert *x509.Certificate
		certFile := fp.Join(mudRootDir, "cert.pem")
		cert, err = pemutil.ReadCertificate(certFile)
		if err != nil {
			return errors.Wrapf(err, "reading certificate from %s failed", certFile)
		}

		// TODO: provide additional information for verification,
		// amongst which are the signature file, the mud file and CA root

		// TODO: write to different location, based on signaturepath and close to the MUD file (if not online)
		newSignaturePath := fp.Join(mudRootDir, "signature.der")
		der, err := ioutil.ReadFile(newSignaturePath)
		if err != nil {
			return errors.Wrap(err, "reading DER signature file failed")
		}

		sd, err := cms.ParseSignedData(der)
		if err != nil {
			return errors.Wrap(err, "parsing signed data failed")
		}

		pool := x509.NewCertPool()
		pool.AddCert(cert)
		options := x509.VerifyOptions{
			Roots: pool, // TODO: make this optional with the CA to trust; now it's like a self-signed cert
		}
		if _, err := sd.VerifyDetached(data, options); err != nil {
			return errors.Wrap(err, "verifying data failed")
		}

		log.Println("MUD verified successfully")

		return nil
	},
}

func getSignatureFilepath(mudfile *mudyang.Mudfile) (string, error) {
	if signatureFlag != "" {
		return signatureFlag, nil
	}
	if mudfile.Mud == nil {
		return "", errors.New("no 'mud' property found in MUD file")
	}
	sig := mudfile.Mud.MudSignature
	if sig == nil {
		return "", errors.New("no signature found in MUD file")
	}
	return *sig, nil
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.PersistentFlags().StringVarP(&signatureFlag, "signature", "s", "", "Location of signature file to use")
}
