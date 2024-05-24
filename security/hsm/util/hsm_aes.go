package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/miekg/pkcs11"
)

func HSM_AES_Key(hsmLib, pin, token, label string) error {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return err
	}

	fmt.Printf("1.  Create AES key, test encryption and decryption.\n")
	// first lookup the key
	aesKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		// pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true), // we do need to extract this
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		// pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, 32), //not permit
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, make([]byte, 32)), /* KeyLength */
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),            /* Name of Key */
		pkcs11.NewAttribute(pkcs11.CKA_ID, label),
	}

	handlers, err := HSM_Find(p, session, aesKeyTemplate, 1)
	if err != nil {
		return err
	}
	if len(handlers) > 0 {
		return fmt.Errorf("FIND PRIVATE KEY")
	}

	aesKey, err := p.CreateObject(session, aesKeyTemplate)
	if err != nil {
		return err
	}
	log.Printf("Created AES Key: %v", aesKey)
	return nil
}

func HSM_AES_CBC_Encrypt(hsmLib, pin, token, label string, data string) ([]byte, error) {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}

	ktemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}

	handlers, err := HSM_Find(p, session, ktemplate, 1)
	if err != nil {
		return nil, err
	}

	if len(handlers) == 0 {
		return nil, fmt.Errorf("NOT FIND PRIVATE KEY")
	}
	aesKey := handlers[0]

	iv := make([]byte, 16)
	_, err = rand.Read(iv)

	if err != nil {
		return nil, err
	}

	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)}
	err = p.EncryptInit(session, mechanism, aesKey)
	if err != nil {
		return nil, err
	}

	ct, err := p.Encrypt(session, []byte(data))
	if err != nil {
		return nil, err
	}

	// append the IV to the ciphertext
	cdWithIV := append(iv, ct...)

	// encryptedData := base64.RawStdEncoding.EncodeToString(cdWithIV)
	// log.Printf("Encrypted IV+Ciphertext %s", encryptedData)

	// // Test data
	// err = p.DecryptInit(session, mechanism, aesKey)
	// if err != nil {
	// 	log.Printf("DecryptInit %s\n", err.Error())

	// 	return nil, err
	// }

	// pt, err := p.Decrypt(session, ct)
	// if err != nil {
	// 	log.Printf("Decrypt %s\n", err.Error())

	// 	return nil, err
	// }
	// log.Printf("Origin %s\n", data)
	// log.Printf("Decrypt %s\n", string(pt))
	return cdWithIV, nil
}

func HSM_AES_CBC_Decrypt(hsmLib string, pin string, token, label string, encrypted []byte) ([]byte, error) {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}

	ktemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}

	handlers, err := HSM_Find(p, session, ktemplate, 1)
	if err != nil {
		return nil, err
	}

	if len(handlers) == 0 {
		return nil, fmt.Errorf("NOT FIND PRIVATE KEY")
	}
	aesKey := handlers[0]

	iv := encrypted[0:16]
	fmt.Printf("iv: %s\n", hex.EncodeToString(iv))

	err = p.DecryptInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)}, aesKey)
	if err != nil {
		fmt.Printf("DecryptInit failed %s\n", err)
		return nil, err
	}
	clear_text, err := p.Decrypt(session, encrypted[16:])
	if err != nil {
		fmt.Printf("Decrypt failed %s\n", err)
		return nil, err
	}
	return clear_text, nil
}

func HSM_AES_Wrapped_By_RSA(hsmLib, pin, token, label string, wrappingLabel string) ([]byte, error) {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}

	privKeyAttrs := []*pkcs11.Attribute{
		// pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}

	aesKeys, err := HSM_Find(p, session, privKeyAttrs, 1)
	if err != nil {
		return nil, err
	}
	if len(aesKeys) == 0 {
		return nil, fmt.Errorf("NOT FOUND PRIVATE KEY")
	}
	aesKey := aesKeys[0]

	pubKeyAttrs := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, wrappingLabel),
	}

	wpbks, err := HSM_Find(p, session, pubKeyAttrs, 1)
	if err != nil {
		return nil, err
	}
	if len(wpbks) == 0 {
		return nil, fmt.Errorf("NOT FOUND WRAPPING KEY")
	}
	wpbk := wpbks[0]

	fmt.Printf("Start wrap key: %d by %d .....\n", aesKey, wpbk)

	// A) wrap AES key using RSA Public Key
	// mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, OAEPSha1Params)}
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_OAEP, OAEPSha1Params)}
	// mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_AES_KEY_WRAP, nil)}

	wrappedPrivBytes, err := p.WrapKey(session, mechanism, wpbk, aesKey)

	// B) wrap AES key using PAD
	// mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_WRAP, nil)}
	// wrappedPrivBytes, err := p.WrapKey(session, mechanism, aesKey, aesKey)
	if err != nil {
		fmt.Printf("WrapKey failed %s \n", err.Error())

		return nil, err
	}
	log.Printf("wrappedPrivBytes: %v", wrappedPrivBytes)
	encryptedPrivKey := base64.RawStdEncoding.EncodeToString(wrappedPrivBytes)
	log.Printf("encryptedPrivKey: %s", encryptedPrivKey)
	return wrappedPrivBytes, nil
}
func HSM_AES_UnWrapped_By_RSA(hsmLib, pin, token, unwrapLabel string, fromLabel string, wrappedPrivBytes []byte) error {

	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return err
	}
	ktemplate1 := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, unwrapLabel),
	}
	privks, err := HSM_Find(p, session, ktemplate1, 1)
	if err != nil {
		return err
	}
	if len(privks) == 0 {
		return fmt.Errorf("NOT FOUND UNWRAP KEY")
	}
	privKey := privks[0]

	aesKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, fromLabel), /* Name of Key */
		pkcs11.NewAttribute(pkcs11.CKA_ID, fromLabel),
	}

	// A) unwrap AES key using RSA key
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_OAEP, OAEPSha1Params)}

	//B) unwrap AES key using PAD
	// mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_WRAP, nil)}

	_, err = p.UnwrapKey(session,
		mechanism,
		privKey,
		wrappedPrivBytes,
		aesKeyTemplate)

	if err != nil {
		return err
	}
	return nil

}
