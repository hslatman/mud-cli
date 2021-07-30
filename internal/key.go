package internal

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.step.sm/crypto/keyutil"
	"go.step.sm/crypto/pemutil"
)

func LoadOrCreateKeyAndChain(chainFilepath, keyFilepath string) ([]*x509.Certificate, crypto.Signer, error) {
	var chain []*x509.Certificate
	var cert *x509.Certificate
	var key crypto.PrivateKey
	var err error
	if !fileExists(keyFilepath) {
		certBytes, keyBytes := generateKey() // TODO: this logic should probably go somewhere different; a key + CSR should be created and sent to CA for a signed cert.
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
		_, err = pemutil.Serialize(key, pemutil.ToFile(keyFilepath, 0600), pemutil.WithPassword([]byte("1234"))) // TODO: provide password or prompt for it
		if err != nil {
			return nil, nil, errors.Wrapf(err, "serializing private key to %s failed", keyFilepath)
		}
	} else {
		chain, err = pemutil.ReadCertificateBundle(chainFilepath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing certificate(s) failed")
		}
		key, err = pemutil.Read(keyFilepath, pemutil.WithPassword([]byte("1234")))
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
