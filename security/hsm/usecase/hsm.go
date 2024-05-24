package usecase

import (
	"encoding/base64"
	"fmt"

	"github.com/blcvn/lib-golang-test/security/hsm/util"
)

func HSM_RSA_Test_Gen() error {
	hsmLib := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	pin := "654321"

	aliceToken := "alice"

	//1. Generate KeyPair
	fmt.Println("--- 1. Generate KeyPair")
	err := util.HSM_RSA_Key(hsmLib, pin, aliceToken, aliceToken)
	if err != nil {
		fmt.Printf("HSM_RSA_Key failed %s\n", err.Error())
		return err
	}

	//2. Get PublicKey
	fmt.Println("--- 2. Get PublicKey")
	pubLabel := fmt.Sprintf("%s_public", aliceToken)
	pKey, err := util.HSM_RSA_Export_PublicKey(hsmLib, pin, aliceToken, pubLabel)
	if err != nil {
		fmt.Printf("HSM_RSA_Export_PublicKey failed %s\n", err.Error())
		return err
	}

	//3. Sign message
	fmt.Println("--- 3. Sign message")
	privLabel := fmt.Sprintf("%s_private", aliceToken)

	message := "THU NGHIEM HE THONG"
	signature, err := util.HSM_RSA_SHA256_Sign(hsmLib, pin, aliceToken, privLabel, message)
	if err != nil {
		fmt.Printf("HSM_RSA_SHA256_Sign failed %s\n", err.Error())
		return err
	}
	fmt.Printf("Signature: %s\n", signature)

	//4. Test signature
	fmt.Println("--- 4. Verify")

	ok := util.RSA_Verify_SHA256(pKey, []byte(message), signature)
	if ok {
		fmt.Println("Verify: OK")
		return nil
	}
	fmt.Println("Verify: Failed")

	return nil
}
func HSM_RSA_Test_Import_PublicKey() error {
	hsmLib := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	pin := "654321"

	message := "THU NGHIEM HE THONG"

	aliceToken := "alice"
	bobToken := "bob"

	privLabel := fmt.Sprintf("%s_private", aliceToken)
	pubLabel := fmt.Sprintf("%s_public", aliceToken)

	aliceLabel := fmt.Sprintf("public_%s", aliceToken)

	// fmt.Println("--- 1. Export Alice Public Key ----")
	pKey, err := util.HSM_RSA_Export_PublicKey(hsmLib, pin, aliceToken, pubLabel)
	if err != nil {
		fmt.Printf("HSM_RSA_Export_PublicKey failed %s\n", err.Error())
		return err
	}

	fmt.Println("--- 2. Import Alice Public Key to Bob HSM ----")
	err = util.HSM_RSA_Import_PublicKey(hsmLib, pin, bobToken, aliceLabel, pKey)
	if err != nil {
		fmt.Printf("HSM_RSA_Import_PublicKey failed %s\n", err.Error())
		return err
	}

	fmt.Println("--- 3. Alice Sign message")

	signature, err := util.HSM_RSA_SHA256_Sign(hsmLib, pin, aliceToken, privLabel, message)
	if err != nil {
		fmt.Printf("HSM_RSA_SHA256_Sign failed %s\n", err.Error())
		return err
	}
	fmt.Printf("Signature: %s\n", signature)

	fmt.Println("--- 4. Bob Verify using HSM ")
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Printf("DecodeString failed %s\n", err.Error())
		return err
	}
	ok, err := util.HSM_RSA_SHA256_Verify(hsmLib, pin, bobToken, aliceLabel, message, sign)
	if err != nil {
		fmt.Printf("HSM_RSA_SHA256_Verify failed %s\n", err.Error())
		return err
	}
	if !ok {
		fmt.Println("Verify: Failed")
		return nil
	}
	fmt.Println("Verify: OK")

	fmt.Println("--- 5. Bob encrypt data with Alice's publicKey ")
	data := []byte("THU NGHIEP RSA ENCRYPTION IN HSM")
	encryptedData, err := util.HSM_RSA_OAEP_Encrypt(hsmLib, pin, bobToken, aliceLabel, data)
	if err != nil {
		fmt.Printf("HSM_RSA_OAEP_Encrypt failed %s\n", err.Error())
		return err
	}

	fmt.Printf("encryptedData: %s \n", base64.StdEncoding.EncodeToString(encryptedData))

	fmt.Println("--- 6. Alice decrypt ciphered data with her privateKey ")
	data1, err := util.HSM_RSA_OAEP_Decrypt(hsmLib, pin, aliceToken, privLabel, encryptedData)
	if err != nil {
		fmt.Printf("HSM_RSA_OAEP_Decrypt failed %s\n", err.Error())
		return err
	}
	fmt.Printf("Origin Data: %s\n", string(data1))
	return nil
}
func HSM_RSA_Test_Wrap_Unwrap() error {
	hsmLib := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	pin := "654321"

	aliceToken := "alice"
	bobToken := "bob"

	// alicePubLabel := fmt.Sprintf("public_%s", aliceToken)
	shareAlice := fmt.Sprintf("shared_%s", aliceToken)
	fromBob := fmt.Sprintf("from_%s", bobToken)

	// fmt.Println("--- 1. Bob Create Alice's shared private Key ")
	// err := util.HSM_AES_Key(hsmLib, pin, bobToken, shareAlice)
	// if err != nil {
	// 	fmt.Printf("HSM_AES_Key failed %s\n", err.Error())
	// 	// return err
	// }

	// fmt.Println("--- 2. Bob Wrap key with Alice's Public Key ")
	// wrappedData, err := util.HSM_AES_Wrapped_By_RSA(hsmLib, pin, bobToken, shareAlice, alicePubLabel)
	// if err != nil {
	// 	fmt.Printf("HSM_AES_Wrapped_By_RSA failed %s\n", err.Error())
	// 	return err
	// }
	// encryptedPrivKey := base64.RawStdEncoding.EncodeToString(wrappedData)
	// log.Printf("encryptedPrivKey: %s", encryptedPrivKey)

	// fmt.Println("--- 3. Alice UnWrap key using her Private's Key ")
	// alicePrivLabel := fmt.Sprintf("%s_private", aliceToken)
	// err = util.HSM_AES_UnWrapped_By_RSA(hsmLib, pin, aliceToken, alicePrivLabel, fromBob, wrappedData)
	// if err != nil {
	// 	fmt.Printf("HSM_AES_UnWrapped_By_RSA failed %s\n", err.Error())
	// 	return err
	// }

	fmt.Println("--- 4. Encrypt Data With new AES Key ")

	message := "FROM ALICE SEND TO BOB"
	cipherData, err := util.HSM_AES_CBC_Encrypt(hsmLib, pin, aliceToken, fromBob, message)
	if err != nil {
		fmt.Printf("HSM_AES_CBC_Encrypt failed: %s\n", err.Error())
		return err
	}
	encryptedData := base64.RawStdEncoding.EncodeToString(cipherData)
	fmt.Printf("EncryptedData: %s\n", encryptedData)

	fmt.Println("--- 5. Test decrypted with old AES Key ")

	originData, err := util.HSM_AES_CBC_Decrypt(hsmLib, pin, bobToken, shareAlice, cipherData)
	if err != nil {
		fmt.Printf("HSM_AES_CBC_Encrypt failed %s\n", err.Error())
		return err
	}
	fmt.Printf("originData: %s\n", originData)

	return nil
}

func HSM_Exchange_Key() error {

	// err = util.HSM_RSA_Key(hsmLib, pin, bobToken, bobToken)
	// if err != nil {
	// 	fmt.Printf("HSM_Gen_RSA failed %s \n", err.Error())
	// 	return err
	// }

	// //2. Exchange public keys
	// fmt.Println("2. Exchange public keys")
	// bobPublicKey, err := util.HSM_RSA_Export_PublicKey(hsmLib, pin, bobToken, bobToken)
	// if err != nil {
	// 	fmt.Printf("HSM_Gen_RSA failed %s \n", err.Error())
	// 	return err
	// }

	// pubKeyPem := util.PublicKeyToBytes(bobPublicKey)

	// fmt.Printf("bobPublicKey: %s", string(pubKeyPem))

	// //3. Generate AES key
	// fmt.Println("3. Generate AES key")
	// err = util.HSM_AES_Key(hsmLib, pin, aliceToken, aliceToken)
	// if err != nil {
	// 	fmt.Printf("HSM_Gen_AES_Key failed %s \n", err.Error())
	// 	return err
	// }

	// //4. Encrypt data with AES key
	// fmt.Println("4. Encrypt data with AES key")
	// data := "THU NGHIEM EXCHANGE KEY"
	// encryptedData, err := util.HSM_AES_CBC_Encrypt(hsmLib, pin, aliceToken, aliceToken, data)
	// if err != nil {
	// 	fmt.Printf("AES_encrypt failed: %s ", err)
	// 	return err
	// }

	// // 5. Encrypt AES key with RSA
	// fmt.Println("5. Encrypt AES key with RSA => Wrap AES key")
	// aesKeyEnc, err := util.HSM_AES_Wrapped_By_RSA(hsmLib, pin, aliceToken, aliceToken)
	// if err != nil {
	// 	fmt.Printf("AES_encrypt failed: %s ", err)
	// 	return err
	// }
	// fmt.Printf("---> aesKeyEnc: %s\n", aesKeyEnc)

	// // 6. Sign message
	// fmt.Println("6. Sign message")
	// aesKeyEncSig, err := util.HSM_RSA_SHA256_Sign(hsmLib, pin, aliceToken, aliceToken, []byte(aesKeyEnc))
	// if err != nil {
	// 	fmt.Printf("RSA_Sign_SHA256 failed: %s ", err)

	// 	return err
	// }

	// //7. Send encrypted data, encrypted AES key, signature
	// fmt.Println("7. Send encrypted data, encrypted AES key, signature")
	// fmt.Printf("---> ecnryptdData: %s \n", encryptedData)
	// fmt.Printf("---> ecnryptdAEKey: %s \n", aesKeyEnc)
	// fmt.Printf("---> aesKeyEncSig: %s \n", aesKeyEncSig)

	// // 8. Verify signature
	// fmt.Println("8. Verify signature")

	// ok, err := util.HSM_RSA_SHA256_Verify(hsmLib, pin, bobToken, bobToken, []byte(aesKeyEnc), aesKeyEncSig)
	// if err != nil {
	// 	fmt.Printf("RSA_Sign_SHA256 failed: %s ", err)

	// 	return err
	// }
	// if !ok {
	// 	fmt.Printf("Encrypted key not correct => Exit")
	// 	return err
	// }
	// fmt.Println("---> Signature ok ")

	// // 9. Decrypt AES key
	// fmt.Println("9. Decrypt AES key (Unwrap)")

	// rsaKey, err := util.HSM_AES_UnWrap_RSA(aesKeyEnc, bobPriv)
	// if err != nil {
	// 	fmt.Printf("DecryptWithPrivateKey failed: %s ", err)
	// 	return err
	// }
	// fmt.Printf("---> rsaKey: %s\n ", rsaKey)

	// // 10. Decrypt data
	// decryptedData, err := util.HSM_AES_Decrypt(hsmLib, pin, bobToken, bobToken, encryptedData)
	// if err != nil {
	// 	fmt.Printf("AES_decrypt failed: %s ", err)
	// 	return err
	// }
	// fmt.Printf("---> decryptedData: %s \n", string(decryptedData))
	return nil
}
