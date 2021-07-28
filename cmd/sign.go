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
	"net/url"
	fp "path/filepath"
	"strings"
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

var baseURLFlag string
var ignoreExistingSignatureFlag bool

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Signs a MUD file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]
		data, err := internal.Read(filepath)
		if err != nil {
			return errors.Wrapf(err, "error reading contents of %s", filepath)
		}
		mudfile, err := internal.Parse(data)
		if err != nil {
			return errors.Wrap(err, "could not get contents")
		}

		// TODO: provide option to first emit standard MUD JSON and use the JSON
		// representation of that as the material to be signed?

		// TODO: look into this logic: if a signature path is know, the signature may already
		// exist or not yet. If it does already exist, we may want to re-sign the MUD file.
		// The location of the new signature can be assumed to be the same in the end.
		// if !ignoreExistingSignatureFlag && mudHasSignature(mudfile) {
		// 	return fmt.Errorf("this MUD already has a signature available at: %s", *mudfile.Mud.MudSignature)
		// }

		existingMudUrl, err := internal.MUDURL(mudfile)
		fmt.Println("existing mud url: ", existingMudUrl)
		if err != nil {
			return errors.Wrap(err, "retrieving MUD URL from MUD failed")
		}

		existingMudSignatureUrl, err := internal.MUDSignatureURL(mudfile)
		fmt.Println("existing mud signature: ", existingMudSignatureUrl)
		if err != nil {
			return errors.Wrap(err, "retrieving MUD signature URL from MUD failed")
		}

		signatureFilename, err := internal.SignatureFilename(filepath)
		fmt.Println("signature path: ", signatureFilename)
		if err != nil {
			return errors.Wrap(err, "retrieving signature path from MUD failed")
		}

		newMudURL := existingMudUrl
		fmt.Println("new MUD url: ", newMudURL)
		newSignatureURL := internal.NewMUDSignatureURL(existingMudUrl, signatureFilename)
		fmt.Println("new signature url: ", newSignatureURL)

		if baseURLFlag != "" {
			newMudURL, err = rewriteBase(newMudURL, baseURLFlag)
			if err != nil {
				return errors.Wrap(err, "rewriting base URL for MUD URL failed")
			}
			newSignatureURL, err = rewriteBase(newSignatureURL, baseURLFlag)
			if err != nil {
				return errors.Wrap(err, "rewriting base URL for MUD signature URL failed")
			}
		}

		fmt.Println("new MUD url: ", newMudURL)
		fmt.Println("new signature url: ", newSignatureURL)

		//var signatureURL *url.URL
		// if baseURLFlag != "" {
		// 	baseURL, err := url.Parse(baseURLFlag)
		// 	if err != nil {
		// 		return errors.Wrap(err, "failed parsing base URL")
		// 	}
		// 	signatureURL = baseURL
		// 	signatureURL.Path = signaturePath // TODO: support path with multiple segments
		// } else {
		// 	signatureURL = mudURL
		// 	signatureURL.Path = signaturePath // TODO: support path with multiple segments
		// }

		// TODO: update Mudfile with location for signature? Needs to be clear that it has indeed be changed.
		// TODO: if signing a local file, provide argument for the full path or directory for the signature file, so
		// that the right value can be added to the MUD file before signing.

		copy, err := ygot.DeepCopy(mudfile)
		if err != nil {
			return errors.Wrap(err, "creating deep copy of MUD YANG representation failed")
		}

		copyMUDFile, ok := copy.(*mudyang.Mudfile)
		if !ok {
			return errors.New("the output MUD YANG is not a *mudyang.Mudfile")
		}

		// TODO: change other properties?
		mudURLString := newMudURL.String()
		copyMUDFile.Mud.MudUrl = &mudURLString
		signatureURLString := newSignatureURL.String()
		copyMUDFile.Mud.MudSignature = &signatureURLString

		diff, err := ygot.Diff(mudfile, copyMUDFile)
		if err != nil {
			return errors.Wrap(err, "diffing the input and output MUD file failed")
		}

		// TODO: can the diff be printed nicely (easily)? It seems to be some text values ...
		log.Println("diff: ", diff)

		differencesFound := len(diff.GetDelete())+len(diff.GetUpdate()) > 0
		if differencesFound {
			// TODO: in case we're using the copy, we probably also need to update the updated_at
			json, err := internal.JSON(copyMUDFile)
			if err != nil {
				return errors.Wrap(err, "getting JSON representation of MUD file failed")
			}
			data = []byte(json)
		}

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

		// TODO: if differences were found, we also should store the new
		// MUD file, so that it can be uploaded later.

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

func rewriteBase(u *url.URL, baseURLString string) (*url.URL, error) {
	// TODO: check this works as expected
	base, err := url.Parse(baseURLString)
	if err != nil {
		return nil, err
	}
	path := u.EscapedPath()
	filename := fp.Base(path)
	baseString := base.String()
	if !strings.HasSuffix(baseString, "/") {
		baseString = baseString + "/"
	}
	newURL, err := url.Parse(baseString + filename)
	if err != nil {
		return nil, err
	}
	newURL.RawQuery = u.RawQuery
	newURL.Fragment = u.Fragment
	return newURL, nil
}

func init() {
	rootCmd.AddCommand(signCmd)

	// TODO: provide a flag that uses a base URL for the MUD URL and signature to exist in the file?

	signCmd.PersistentFlags().StringVarP(&baseURLFlag, "base-url", "u", "", "Base URL to use for MUD URL and signature location")
	signCmd.PersistentFlags().StringVarP(&signatureFlag, "signature", "s", "", "Location of signature file to set")
	signCmd.PersistentFlags().BoolVar(&ignoreExistingSignatureFlag, "ignore-existing-signature", false, "Ignore")
}
