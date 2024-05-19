package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"

	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/cloudflare/cfssl/log"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Error(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Error(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		fmt.Printf("BytesToPrivateKey: is encrypted pem block \n")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Printf("BytesToPrivateKey: DecryptPEMBlock failed %s \n", err.Error())
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		fmt.Printf("BytesToPrivateKey: ParsePKCS1PrivateKey failed %s \n", err.Error())
		return nil, err
	}
	return key, nil
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey_PKCS8(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		fmt.Printf("BytesToPrivateKey: is encrypted pem block \n")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Printf("BytesToPrivateKey: DecryptPEMBlock failed %s \n", err.Error())
			return nil, err
		}
	}
	key, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		fmt.Printf("BytesToPrivateKey: ParsePKCS1PrivateKey failed %s \n", err.Error())
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		fmt.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Error(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Error(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Error("not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) (string, error) {
	// hash := sha512.New()
	hash := sha1.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// rsa_public_key = RSA.importKey(publicKey)
//     rsa_public_key = PKCS1_OAEP.new(rsa_public_key)
//     encrypted_text = rsa_public_key.encrypt(passphrase)
//     encrypted_base64 = base64.b64encode(encrypted_text)

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext string, priv *rsa.PrivateKey) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	// hash := sha512.New()
	hash := sha1.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, sDec, nil)
	return string(plaintext), err
}

// encrypted_text = base64.b64decode(encrypted_base64)
// private_key = RSA.importKey(privateKey)
// cipher = PKCS1_OAEP.new(private_key)
// decrypted_text = cipher.decrypt(encrypted_text)

func LoadPrivateKeyFile(file string) (*rsa.PrivateKey, error) {
	privBytes, err := ioutil.ReadFile(file) // This is fine with Encryption
	if err != nil {
		fmt.Printf("LoadKeyFile: open file failed %s\n", err)
		return nil, err
	}
	return BytesToPrivateKey(privBytes)
}
