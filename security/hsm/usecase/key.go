package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/blcvn/lib-golang-test/security/hsm/util"
)

func HSM_Exchange_AES_Key() {
	//1. Generate AES Key and Encrypt Data
	hsmLib := "/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
	pin := ""
	walletToken := "wallet"

	fmt.Printf("1. Key_Gen_AES_Key \n")
	util.Key_Gen_AES_Key(hsmLib, pin, walletToken, walletToken)

	msg := "Thu nghiem ma hoa"

	fmt.Printf("2. Key_Encrypt_With_AES_Key \n")
	util.Key_Encrypt_With_AES_Key(hsmLib, pin, "wallet", "wallet-aes", msg)
	//2. Gen RSA Key => Get public Cert

	//3. Get Wrapped Key

	//4. Import wrapped key to decrypt data

	// Init PKCS
	//  Create RSA Key used for wrapped transfer

	// // B) unwrap AES key using RSA Public Wrapping Key
	// ik, err := p.UnwrapKey(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)}, wpvk, wrappedPrivBytes, aesKeyTemplate)

	// if err != nil {
	// 	log.Fatalf("Unwrap Failed: %v", err)
	// }

	// // use unwraped key to decrypt the same data we did at the beginning
	// err = p.DecryptInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, cdWithIV[0:16])}, ik)
	// if err != nil {
	// 	panic(fmt.Sprintf("EncryptInit() failed %s\n", err))
	// }

	// pt, err = p.Decrypt(session, ct[:16])
	// if err != nil {
	// 	panic(fmt.Sprintf("Encrypt() failed %s\n", err))
	// }

	// log.Printf("Decrypt %s", string(pt))

}
func EncryptKey() {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	plaintext := []byte("exampleplaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)
}

func Decode() {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
	nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", plaintext)
}
