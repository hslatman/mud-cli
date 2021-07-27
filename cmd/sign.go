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
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	fp "path/filepath"
	"time"

	cms "github.com/github/ietf-cms"
	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/openconfig/ygot/ygot"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.step.sm/crypto/keyutil"
	"go.step.sm/crypto/pemutil"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
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
			return errors.Wrapf(err, "retrieving contents from %s failed", filepath)
		}

		mudfile := &mudyang.Mudfile{}
		if err := mudyang.Unmarshal(json, mudfile); err != nil {
			return errors.Wrapf(err, "unmarshaling JSON failed")
		}

		// TODO: provide flag to override an existing signature
		// if mudHasSignature(mud) {
		// 	fmt.Printf("this MUD already has a signature available at: %s\n", *mud.Mud.MudSignature)
		// 	return
		// }

		signaturePath, err := internal.SignaturePath(filepath)
		fmt.Println("signature path: ", signaturePath)
		if err != nil {
			return errors.Wrap(err, "retrieving signature path from MUD failed")
		}

		// TODO: update Mudfile with location for signature? Needs to be clear that it has indeed be changed.
		// TODO: if signing a local file, provide argument for the full path or directory for the signature file, so
		// that the right value can be added to the MUD file before signing.

		output, err := ygot.DeepCopy(mudfile)
		if err != nil {
			return errors.Wrap(err, "creating deep copy of MUD YANG representation failed")
		}

		outputMudfile, ok := output.(*mudyang.Mudfile)
		if !ok {
			return errors.New("the output MUD YANG is not a *mudyang.Mudfile")
		}

		// TODO: if changing the signature, should we update the updated_at too? And/or others?
		outputMudfile.Mud.MudSignature = &signaturePath

		n, err := ygot.Diff(mudfile, outputMudfile)
		if err != nil {
			return errors.Wrap(err, "diffing the input and output MUD file failed")
		}

		// TODO: after diffing, check only expected number of changes were made
		log.Println("diff: ", n)

		//jsonString, err := ygot.EmitJSON(outputMudfile, &ygot.EmitJSONConfig{
		jsonString, err := ygot.EmitJSON(mudfile, &ygot.EmitJSONConfig{
			Format: ygot.RFC7951,
			Indent: "  ",
			RFC7951Config: &ygot.RFC7951JSONConfig{
				AppendModuleName: true,
			},
			SkipValidation: false,
			// TODO: other validation options?
		})
		if err != nil {
			return errors.Wrap(err, "could not marshal MUD file into JSON")
		}

		data := []byte(jsonString)

		fmt.Println(data)

		// TODO: allow to provide a cert and key (and/or fallback to some stored in .mud dir?)
		var cert *x509.Certificate
		var key crypto.PrivateKey
		keyFile := fp.Join(mudRootDir, "key.pem")
		certFile := fp.Join(mudRootDir, "cert.pem")
		if !fileExists(keyFile) {
			certBytes, keyBytes := generateKey() // TODO: this logic should probably go somewhere different; a key + CSR should be created and sent to CA for a signed cert.
			cert, err = x509.ParseCertificate(certBytes)
			if err != nil {
				return errors.Wrap(err, "parsing certificate failed")
			}
			_, err = pemutil.Serialize(cert, pemutil.ToFile(certFile, 0600))
			if err != nil {
				return errors.Wrapf(err, "serializing certificate to %s failed", certFile)
			}
			key, err = x509.ParsePKCS8PrivateKey(keyBytes)
			if err != nil {
				return errors.Wrap(err, "parsing private key failed")
			}
			_, err = pemutil.Serialize(key, pemutil.ToFile(keyFile, 0600), pemutil.WithPassword([]byte("1234"))) // TODO: provide password or prompt for it
			if err != nil {
				return errors.Wrapf(err, "serializing private key to %s failed", keyFile)
			}
		} else {
			cert, err = pemutil.ReadCertificate(certFile)
			if err != nil {
				return errors.Wrapf(err, "reading certificate from %s failed", certFile)
			}
			k, err := pemutil.Read(keyFile, pemutil.WithPassword([]byte("1234")))
			if err != nil {
				return errors.Wrapf(err, "reading private key from %s failed", keyFile)
			}
			key = k
		}

		fmt.Println(cert)
		fmt.Println(key)
		fmt.Println(fmt.Sprintf("%T", key))

		// TODO: add intermediates? Or how to integrate with a (private, trusted) CA? RFC says that they MUST be added.
		certs := []*x509.Certificate{cert}

		// TODO: allow signing with some other signer, too?
		signer, ok := key.(crypto.Signer)
		if !ok {
			return errors.New("key is not a signer")
		}

		// TODO: add signed timestamp?
		signature, err := cms.SignDetached(data, certs, signer)
		if err != nil {
			return errors.Wrap(err, "signing data failed")
		}

		// TODO: write to different location, based on signaturepath and close to the MUD file
		// TODO: also provide an option to encode it to PEM instead?
		newSignaturePath := fp.Join(mudRootDir, "signature.der")
		err = ioutil.WriteFile(newSignaturePath, signature, 0644)
		if err != nil {
			return errors.Wrap(err, "writing DER signature failed")
		}

		fmt.Println(signature)

		// TODO: print output on how/where to store the file + signature?

		log.Println("MUD file signed successfully")
		return nil
	},
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func generateKey() (certBytes, keyBytes []byte) {
	priv, err := keyutil.GenerateKey("EC", "P-256", 0)
	//priv, err := keyutil.GenerateKey("RSA", "", 2048)
	if err != nil {
		log.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"mudsign example organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature, // TODO: fix usage? Would actually be better if this was more like a separate command, I guess
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("failed to create certificate: %s", err)
	}

	keyBytes, err = x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("failed to marshal private key: %s", err)
	}

	return derBytes, keyBytes
}

func mudHasSignature(mud *mudyang.Mudfile) bool {
	return mud.Mud.MudSignature != nil
}

func init() {
	rootCmd.AddCommand(signCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
