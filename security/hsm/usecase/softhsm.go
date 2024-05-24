package usecase

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/blcvn/lib-golang-test/security/hsm/util"
)

func TestImportPrivateKey() {
	lib_path := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	token := "test_import"
	pin := "8764329"
	key_file := "test/hsm/cert/priv_sk"

	privateKey, err := util.ECDSA_LoadPrivateKeyFile(key_file)
	if err != nil {
		fmt.Println("ECDSA_LoadPrivateKeyFile error: ", err)
		return
	}
	err = util.HSM_ImportKey(lib_path, pin, token, privateKey)
	if err != nil {
		fmt.Println("HSM_ImportKey error: ", err)
		return
	}
	fmt.Println("OK")

}
func TestSignAndVerify() {
	lib_path := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	token := "test_orderer"
	pin := "8764329"
	msg := "THu nghiem sign"

	key_file := "test/hsm/cert/priv_sk"
	privateKey, err := util.ECDSA_LoadPrivateKeyFile(key_file)
	if err != nil {
		fmt.Println("ECDSA_LoadPrivateKeyFile error: ", err)
		return
	}
	cert_file := "test/hsm/cert/orderer0.test.shard1.com-cert.pem"
	publicKey, err := util.ECDSA_LoadPublicKeyFile(cert_file)
	if err != nil {
		fmt.Println("ECDSA_LoadPublicKeyFile error: ", err)
		return
	}
	hash := sha256.Sum256([]byte(msg))

	fmt.Println("Test signature from file ......")
	sig1, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		fmt.Printf("ecdsa.SignASN1 failed %s ", err.Error())
		return
	}
	// var esig1 struct {
	// 	R, S *big.Int
	// }
	// if _, err := asn1.Unmarshal(sig1, &esig1); err != nil {
	// 	fmt.Println("asn1.Unmarshal error: ", err)
	// 	return
	// }
	// ok1 := ecdsa.Verify(publicKey, hash[:], esig1.R, esig1.S)
	// if !ok1 {
	// 	fmt.Println("Test2 failed")
	// 	return
	// }
	ok1 := ecdsa.VerifyASN1(publicKey, hash[:], sig1)
	if !ok1 {
		fmt.Println("ok1 failed")
		return
	}
	fmt.Println("OK1")

	fmt.Println("Test signature from HSM ......")
	sig2, err := util.HSM_RSA_SHA256_Sign(lib_path, pin, token, "test_import", msg)
	if err != nil {
		fmt.Println("HSM_Sign error: ", err)
		return
	}
	fmt.Printf("Signature: %s\n", sig2)

	// var esig2 struct {
	// 	R, S *big.Int
	// }
	// if _, err := asn1.Unmarshal(sig2, &esig2); err != nil {
	// 	fmt.Println("asn1.Unmarshal error: ", err)
	// 	return
	// }
	// ok2 := ecdsa.Verify(publicKey, hash[:], esig2.R, esig2.S)
	// if !ok2 {
	// 	fmt.Println("Test2 failed")
	// 	return
	// }

	sig2Data, err := base64.StdEncoding.DecodeString(sig2)
	if err != nil {
		fmt.Println("DecodeString error: ", err)
		return
	}
	ok := ecdsa.VerifyASN1(publicKey, hash[:], sig2Data)
	if !ok {
		fmt.Println("failed")
		return
	}
	fmt.Println("Test2 OK")

}
func ReadObject() {
	lib_path := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	token := "test"
	pin := "8764329"
	util.HSM_ReadObject(lib_path, token, pin)
}

func GenSKI() {
	privateFile := "test/hsm/cert/priv_sk"
	pubKey, err := util.ECDSA_LoadPrivateKeyFile(privateFile)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
		return
	}
	ski := util.GenSKI(&pubKey.PublicKey)
	fmt.Printf("SKI: %s\n", ski)

	// certFile := "test/hsm/cert/orderer0.test.shard1.com-cert.pem"
	// str, err := util.GetCertSKI(certFile)
	// if err != nil {
	// 	fmt.Printf("Error: %s \n", err.Error())
	// 	return
	// }
	// fmt.Printf("result: %s\n", str)
	// pubKey, err := util.ECDSA_LoadPublicKeyFile(certFile)
	// if err != nil {
	// 	fmt.Printf("Error: %s \n", err.Error())
	// 	return
	// }
	// ski := util.GenSKI(pubKey)
	// fmt.Printf("SKI: %s\n", ski)

}
