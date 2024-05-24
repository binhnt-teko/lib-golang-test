package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/blcvn/lib-golang-test/security/hsm/util"
)

func Exchange_AES_Key() error {
	fmt.Println("1. Generate RSA key pairs")

	//1. Generate RSA key pairs
	alicePriv, alicePub := util.GenerateKeyPair(256 * 8)
	bobPriv, bobPub := util.GenerateKeyPair(256 * 8)

	fmt.Println("2. Exchange public keys")
	//2. Exchange public keys
	fmt.Printf("---> Alice Public Key: \n %s \n ", util.PublicKeyToBytes(alicePub))
	fmt.Printf("---> Bob Public Key: \n  %s \n ", util.PublicKeyToBytes(bobPub))

	fmt.Println("3. Generate AES key")
	//3. Generate AES key
	randData := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randData); err != nil {
		fmt.Printf("ReadFull failed: %s ", err)
		return err
	}
	AES_KEY := hex.EncodeToString(randData)
	fmt.Printf("---> AES_KEY: %s\n", AES_KEY)

	fmt.Println("4. Encrypt data with AES key")

	//4. Encrypt data with AES key
	data := "THU NGHIEM EXCHANGE KEY"

	encryptedData, err := util.AES_encrypt(data, AES_KEY)
	if err != nil {
		fmt.Printf("AES_encrypt failed: %s ", err)

		return err
	}
	fmt.Printf("---> encryptedData: %s\n", encryptedData)

	fmt.Println("5. Encrypt AES key with RSA")

	// 5. Encrypt AES key with RSA
	fmt.Printf("---> Public key size: %d \n", bobPub.Size())

	aesKeyEnc, err := util.EncryptWithPublicKey([]byte(AES_KEY), bobPub)
	if err != nil {
		fmt.Printf("EncryptWithPublicKey failed: %s ", err)

		return err
	}
	fmt.Printf("---> aesKeyEnc: %s\n", aesKeyEnc)

	// 6. Sign message
	fmt.Println("6. Sign message")

	aesKeyEncSig, err := util.RSA_Sign_SHA256([]byte(aesKeyEnc), alicePriv)
	if err != nil {
		fmt.Printf("RSA_Sign_SHA256 failed: %s ", err)

		return err
	}
	fmt.Printf("---> aesKeyEncSig: %s\n", aesKeyEncSig)

	//7. Send encrypted data, encrypted AES key, signature
	fmt.Println("7. Send encrypted data, encrypted AES key, signature")

	fmt.Printf("---> ecnryptdData: %s \n", encryptedData)
	fmt.Printf("---> ecnryptdAEKey: %s \n", aesKeyEnc)
	fmt.Printf("---> aesKeyEncSig: %s \n", aesKeyEncSig)

	// 8. Verify signature
	fmt.Println("8. Verify signature")

	ok := util.RSA_Verify_SHA256(alicePub, []byte(aesKeyEnc), aesKeyEncSig)
	if !ok {
		fmt.Printf("Encrypted key not correct => Exit")
		return fmt.Errorf("VERIFIED FAILURE")
	}
	fmt.Println("---> Signature ok ")

	// 9. Decrypt AES key
	fmt.Println("9. Decrypt AES key")

	rsaKey, err := util.DecryptWithPrivateKey(aesKeyEnc, bobPriv)
	if err != nil {
		fmt.Printf("DecryptWithPrivateKey failed: %s ", err)
		return err
	}
	fmt.Printf("---> rsaKey: %s\n ", rsaKey)

	fmt.Println("10. Decrypt data")

	//10. Decrypt data
	decryptedData, err := util.AES_decrypt(encryptedData, rsaKey)
	if err != nil {
		fmt.Printf("AES_decrypt failed: %s ", err)
		return err
	}
	fmt.Printf("---> decryptedData: %s \n", string(decryptedData))
	return nil
}
