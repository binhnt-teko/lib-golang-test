package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func ECDSA_LoadPrivateKeyFile(privateKeyFile string) (*ecdsa.PrivateKey, error) {
	certBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	privBlock, _ := pem.Decode(certBytes)
	if privBlock == nil {
		return nil, fmt.Errorf("PEM DECODE FAILED")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	return privKey.(*ecdsa.PrivateKey), err
}
func ECDSA_LoadPublicKeyFile(certFile string) (*ecdsa.PublicKey, error) {
	certBytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil {
		return nil, fmt.Errorf("PEM DECODE FAILED")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)

	if err != nil {
		fmt.Println("parsex509:", err.Error())
		return nil, err
	}
	return cert.PublicKey.(*ecdsa.PublicKey), nil
}

func GetCertSKI(certFile string) (string, error) {
	certBytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		return "", err
	}
	certBlock, _ := pem.Decode(certBytes)

	cert, err := x509.ParseCertificate(certBlock.Bytes)

	if err != nil {
		fmt.Println("parsex509:", err.Error())
		return "", err
	}
	ret := hex.EncodeToString(cert.SubjectKeyId)
	fmt.Println(ret)
	return ret, nil

}
func GenSKI(pubKey *ecdsa.PublicKey) string {
	// Marshall the public key
	raw := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	ski := hash.Sum(nil)
	ret := hex.EncodeToString(ski)
	fmt.Println(ret)
	return ret
}
