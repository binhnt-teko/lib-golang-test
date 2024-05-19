package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/blcvn/lib-golang-test/security/hsm/util"
)

func TestAES_CBC() {
	lib_path := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	token := "test1"
	pin := "8764329"

	err := util.HSM_AES_Key(lib_path, pin, token)
	// if err != nil {
	// 	fmt.Printf("HSM_AES_Key: error: %s\n ", err)
	// 	return
	// }
	msg := "thu nghiem encrypt "
	encrypted, err := util.HSM_AES_CBC_Encrypt(lib_path, pin, token, []byte(msg))
	if err != nil {
		fmt.Printf("HSM_AES_CBC_Encrypt: error: %s\n", err)
		return
	}
	fmt.Printf("Encrypted: %s\n", base64.RawStdEncoding.EncodeToString(encrypted))
	clear_text, err := util.HSM_AES_CBC_Decrypt(lib_path, pin, token, encrypted)
	if err != nil {
		fmt.Printf("HSM_AES_CBC_Decrypt: error: %s \n", err)
		return
	}
	fmt.Printf("Text: %s", string(clear_text))

}

func AES_CFB() {
	key := []byte("0123456789012345")
	iv := []byte("9876543210987654")
	msg := []byte("Hello AES, my old friend")

	enc, err := util.AES_Encrypt(msg, key, iv)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", enc)
	cryptoText := "RNaMzLQHPtdlk+tIiJUrRptdBHDikrjIfrMS3ywS5HRe1EJ+N7bte0kyKQ=="
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv = ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	fmt.Printf("%s", ciphertext)
}
