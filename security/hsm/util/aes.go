package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const BLOCK_SIZE = 16

func trans(key string) []byte {
	kb := md5.Sum([]byte(key))
	bs := kb[:]
	return bs
}

func AES_encrypt(stringToEncrypt string, passphrase string) (string, error) {
	key := trans(passphrase)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AES_decrypt(cryptoText string, passphrase string) ([]byte, error) {
	key := trans(passphrase)

	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	plaintext := make([]byte, len(ciphertext))

	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// Encrypt encrypts the plaintext with key and iv using AES-CFB and returns it in base64.
func AES_Encrypt(plaintext, key, iv []byte) ([]byte, error) {
	// Create ciphertext.
	ciphertext := make([]byte, len(plaintext))
	// Create AES cipher.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Get AES-CFB stream encrypted.
	stream := cipher.NewCFBEncrypter(aesBlock, iv)

	// Encrypt the msg and store the results in ciphertext.
	stream.XORKeyStream(ciphertext, plaintext)
	// Base64 encode it.
	encodedCiphertext := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(encodedCiphertext, ciphertext)
	return encodedCiphertext, nil
}
