package internal

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/smallstep/cli/ui"
	"go.step.sm/crypto/keyutil"
	"go.step.sm/crypto/pemutil"
)

func LoadOrCreateKeyAndChain(chainFilepath, keyFilepath string) ([]*x509.Certificate, crypto.Signer, error) {
	var chain []*x509.Certificate
	var cert *x509.Certificate
	var key crypto.PrivateKey
	var err error
	if !fileExists(keyFilepath) {
		// TODO: split logic for the key and chain/cert? Or make clear that this is for testing/demo purposes?
		shouldContinue, err := ui.PromptYesNo(
			fmt.Sprintf("key at %s does not exist; create a new one?", keyFilepath),
			ui.WithRichPrompt(),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error prompting user")
		}
		if !shouldContinue {
			return nil, nil, errors.New("no private key available nor created")
		}
		certBytes, keyBytes, err := generateKey()
		if err != nil {
			return nil, nil, errors.Wrap(err, "error generating new private key")
		}
		cert, err = x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing certificate failed")
		}
		_, err = pemutil.Serialize(cert, pemutil.ToFile(chainFilepath, 0600))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "serializing certificate to %s failed", chainFilepath)
		}
		chain = []*x509.Certificate{cert}
		key, err = x509.ParsePKCS8PrivateKey(keyBytes)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing private key failed")
		}
		options := []pemutil.Options{}
		options = append(options, pemutil.ToFile(keyFilepath, 0600))
		password, err := ui.PromptPasswordGenerate(
			"Please enter a password for the private key [a random password will be generated if you leave this empty]",
			ui.WithRichPrompt(),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error prompting user for password")
		}
		options = append(options, pemutil.WithPassword(password))
		_, err = pemutil.Serialize(key, options...)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "serializing private key to %s failed", keyFilepath)
		}
	} else {
		chain, err = pemutil.ReadCertificateBundle(chainFilepath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing certificate(s) failed")
		}
		options := []pemutil.Options{}
		options = append(options, pemutil.WithPasswordPrompt(
			fmt.Sprintf("Please enter the password to decrypt %s", keyFilepath),
			func(s string) ([]byte, error) {
				return ui.PromptPassword(s)
			}),
		)
		key, err = pemutil.Read(keyFilepath, options...)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "reading private key from %s failed", keyFilepath)
		}
	}
	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, nil, errors.New("key is not a signer")
	}
	return chain, signer, nil
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func generateKey() (certBytes, keyBytes []byte, err error) {
	private, err := keyutil.GenerateKey("EC", "P-256", 0)
	if err != nil {
		return certBytes, keyBytes, err
	}
	public, err := keyutil.PublicKey(private)
	if err != nil {
		return certBytes, keyBytes, err
	}
	i, err := rand.Int(rand.Reader, big.NewInt(100000000000000000))
	if err != nil {
		return certBytes, keyBytes, err
	}
	template := x509.Certificate{
		SerialNumber: i,
		Subject: pkix.Name{
			Organization: []string{"mud-cli example signing organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature, // TODO: fix usage? Would actually be better if this was more like a separate command, I guess
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	certBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, public, private)
	if err != nil {
		return certBytes, keyBytes, err
	}
	keyBytes, err = x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return certBytes, keyBytes, err
	}
	return certBytes, keyBytes, nil
}
