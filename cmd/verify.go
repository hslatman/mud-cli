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
	"crypto/x509"
	"fmt"
	"log"

	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/pkg/errors"
	"github.com/smallstep/pkcs7"
	"github.com/spf13/cobra"
	"go.step.sm/crypto/pemutil"
)

var signatureFlag string
var caBundleFilepathFlag string

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
		signaturePath, err := getSignatureFilepath(mudfile)
		if err != nil {
			return errors.Wrap(err, "retrieving signature from MUD failed")
		}

		var caBundle []*x509.Certificate
		if caBundleFilepathFlag != "" {
			caBundle, err = pemutil.ReadCertificateBundle(caBundleFilepathFlag)
			if err != nil {
				return errors.Wrapf(err, "reading certificate from %s failed", caBundleFilepathFlag)
			}
		}

		der, err := internal.Read(signaturePath)
		if err != nil {
			return errors.Wrap(err, "reading DER signature failed")
		}

		p7, err := pkcs7.Parse(der)
		if err != nil {
			return fmt.Errorf("failed parsing signed data: %w", err)
		}

		// add the content, as it was detached before
		p7.Content = data

		var roots *x509.CertPool
		if len(caBundle) > 0 {
			roots = x509.NewCertPool()
			for _, cert := range caBundle {
				roots.AddCert(cert)
			}
		}

		// // TODO: more verify options? and stricter?
		// options := x509.VerifyOptions{
		// 	//CurrentTime: time.Date(2021, 7, 16, 10, 1, 1, 0, time.Local), // TODO: remove this fully; or make it some kind of option to skip the time check on the cert(s)?
		// 	Roots: roots,
		// }

		if err := p7.VerifyWithChain(roots); err != nil {
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
	verifyCmd.PersistentFlags().StringVarP(&caBundleFilepathFlag, "ca-bundle", "c", "", "Path to CA (root) certificates to trust (PEM)")
}
