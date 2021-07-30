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
	"io/ioutil"
	"log"
	"net/url"
	"os"
	fp "path/filepath"
	"strings"
	"time"

	cms "github.com/github/ietf-cms"
	"github.com/hslatman/mud-cli/internal"
	"github.com/hslatman/mud.yang.go/pkg/mudyang"
	"github.com/openconfig/ygot/ygot"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var baseURLFlag string
var ignoreExistingSignatureFlag bool
var normalizeFlag bool
var keyFilepathFlag string
var chainFilepathFlag string

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
		if normalizeFlag {
			json, err := internal.JSON(mudfile)
			if err != nil {
				return errors.Wrap(err, "normalizing MUD file failed")
			}
			data = []byte(json)
		}

		existingMudUrl, err := internal.MUDURL(mudfile)
		if err != nil {
			return errors.Wrap(err, "retrieving MUD URL from MUD failed")
		}

		// TODO: look into this logic: if a signature path is know, the signature may already
		// exist or not yet. If it does already exist, we may want to re-sign the MUD file.
		// The location of the new signature can be assumed to be the same in the end.
		// if !ignoreExistingSignatureFlag && mudHasSignature(mudfile) {
		// 	return fmt.Errorf("this MUD already has a signature available at: %s", *mudfile.Mud.MudSignature)
		// }
		// existingMudSignatureUrl, err := internal.MUDSignatureURL(mudfile)
		// if err != nil {
		// 	return errors.Wrap(err, "retrieving MUD signature URL from MUD failed")
		// }

		mudFilename, err := internal.MUDFilename(filepath)
		if err != nil {
			return errors.Wrap(err, "retrieving MUD filename failed")
		}

		signatureFilename, err := internal.SignatureFilename(filepath)
		if err != nil {
			return errors.Wrap(err, "retrieving signature path from MUD failed")
		}

		newMudURL := existingMudUrl
		newSignatureURL := internal.NewMUDSignatureURL(existingMudUrl, signatureFilename)

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

		copy, err := ygot.DeepCopy(mudfile)
		if err != nil {
			return errors.Wrap(err, "creating deep copy of MUD YANG representation failed")
		}

		copyMUDFile, ok := copy.(*mudyang.Mudfile)
		if !ok {
			return errors.New("the output MUD YANG is not a *mudyang.Mudfile")
		}

		// TODO: change more properties?
		mudURLString := newMudURL.String()
		copyMUDFile.Mud.MudUrl = &mudURLString
		signatureURLString := newSignatureURL.String()
		copyMUDFile.Mud.MudSignature = &signatureURLString

		diff, err := ygot.Diff(mudfile, copyMUDFile)
		if err != nil {
			return errors.Wrap(err, "diffing the input and output MUD file failed")
		}

		// TODO: can the diff be printed nicely (easily)? It seems to be some text values ...
		//log.Println("diff: ", diff)

		mudHadDifferences := len(diff.GetDelete())+len(diff.GetUpdate()) > 0 || normalizeFlag
		if mudHadDifferences {
			now := time.Now().Format("2006-01-02T15:04:05Z07:00")
			copyMUDFile.Mud.LastUpdate = &now
			json, err := internal.JSON(copyMUDFile)
			if err != nil {
				return errors.Wrap(err, "getting JSON representation of MUD file failed")
			}
			data = []byte(json)
		}

		chain, signer, err := internal.LoadOrCreateKeyAndChain(chainFilepathFlag, keyFilepathFlag)
		if err != nil {
			return errors.Wrap(err, "loading/creating private key failed")
		}

		// TODO: prevent signing with certificate that is no longer valid (or almost going to expire?)

		// TODO: allow signing with some other signer, not based on key and cert file, too?
		signature, err := cms.SignDetached(data, chain, signer)
		if err != nil {
			return errors.Wrap(err, "signing data failed")
		}

		// TODO: if MUD file is local file (not URL), the put it next to the MUD file?
		outputDir := fp.Join(mudRootDir, "files")
		if !dirExists(outputDir) {
			os.MkdirAll(outputDir, 0700)
		}

		// TODO: also provide an option to encode it to PEM instead? Have only seend DER examples, though.
		newSignatureFilepath := fp.Join(outputDir, signatureFilename)
		err = ioutil.WriteFile(newSignatureFilepath, signature, 0644)
		if err != nil {
			return errors.Wrap(err, "writing DER signature failed")
		}

		log.Printf("MUD signature successfully written to %s\n", newSignatureFilepath)

		if mudHadDifferences {
			newMUDFilepath := fp.Join(outputDir, mudFilename)
			err = ioutil.WriteFile(newMUDFilepath, data, 0644)
			if err != nil {
				return errors.Wrap(err, "writing DER signature failed")
			}
			log.Printf("Updated MUD file written to %s\n", newMUDFilepath)
		}

		return nil
	},
}

// func mudHasSignature(mud *mudyang.Mudfile) bool {
// 	return mud.Mud.MudSignature != nil
// }

func rewriteBase(u *url.URL, baseURLString string) (*url.URL, error) {
	// TODO: check/test this works as expected in all/most cases
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

	defaultKeyFilepath := fp.Join(mudRootDir, "key.pem")
	defaultChainFilepath := fp.Join(mudRootDir, "chain.pem")

	signCmd.PersistentFlags().StringVarP(&baseURLFlag, "base-url", "u", "", "Base URL to use for MUD URL and signature location")
	signCmd.PersistentFlags().StringVarP(&signatureFlag, "signature", "s", "", "Location of signature file to set")
	signCmd.PersistentFlags().BoolVar(&ignoreExistingSignatureFlag, "ignore-existing-signature", false, "Ignore case in which MUD already has a signature")
	signCmd.PersistentFlags().BoolVar(&normalizeFlag, "normalize", false, "Normalize the MUD JSON according to default mud.yang.go order")
	signCmd.PersistentFlags().StringVarP(&keyFilepathFlag, "key", "k", defaultKeyFilepath, "Path to private key file (PEM)")
	signCmd.PersistentFlags().StringVarP(&chainFilepathFlag, "chain", "c", defaultChainFilepath, "Path to certificate chain or self-signed certificate (PEM)")
}
