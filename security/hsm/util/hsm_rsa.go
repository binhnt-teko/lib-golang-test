package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/cloudflare/cfssl/log"
	"github.com/miekg/pkcs11"
)

func HSM_RSA_Key(hsmLib, pin, token, label string) error {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return err
	}

	// first lookup the key
	pubLabel := fmt.Sprintf("%s_public", label)
	wpublicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		// pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, pubLabel),
		pkcs11.NewAttribute(pkcs11.CKA_ID, pubLabel),
	}
	privLabel := fmt.Sprintf("%s_private", label)

	wprivateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		// pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, privLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
		pkcs11.NewAttribute(pkcs11.CKA_ID, privLabel),
	}
	_, _, err = p.GenerateKeyPair(session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)},
		wpublicKeyTemplate, wprivateKeyTemplate)

	if err != nil {
		log.Fatalf("failed to generate keypair: %s\n", err)
		return err
	}
	return nil
}
func HSM_RSA_OAEP_Encrypt(libPath string, pin string, token string, label string, data []byte) ([]byte, error) {
	//1. Init HSM
	p, cancel, session, err := HSM_Login(libPath, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}
	handles, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return nil, err
	}
	if len(handles) == 0 {
		return nil, fmt.Errorf("NOT FOUND KEY")
	}

	//Get first key, RSA_PKCS_OAEP with SHA1
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_OAEP, OAEPSha1Params)}
	err = p.EncryptInit(session, mechanism, handles[0])
	if err != nil {
		fmt.Printf("EncryptInit error: %s ", err.Error())
		return nil, err
	}

	encryptedData, err := p.Encrypt(session, []byte(data))
	if err != nil {
		fmt.Println("DecryptFinal error: ", err)
		return nil, nil
	}

	return encryptedData, nil
}

func HSM_RSA_OAEP_Decrypt(libPath string, pin string, token string, label string, sDec []byte) ([]byte, error) {
	//1. Init HSM
	p, cancel, session, err := HSM_Login(libPath, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}
	handles, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return nil, err
	}
	//Get first key, RSA_PKCS_OAEP with SHA1
	err = p.DecryptInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_OAEP, OAEPSha1Params)}, handles[0])
	if err != nil {
		fmt.Printf("DecryptInit error: %s ", err.Error())
		return nil, err
	}

	clear_text, err := p.Decrypt(session, []byte(sDec))
	// clear_text, err := p.DecryptFinal(session)
	if err != nil {
		fmt.Println("DecryptFinal error: ", err)
		return nil, nil
	}
	return clear_text, nil
}

func HSM_RSA_SHA256_Sign(libPath string, pin string, token string, label string, msg string) (string, error) {
	//1. Init HSM
	p, cancel, session, err := HSM_Login(libPath, pin, token)
	defer cancel(session)
	if err != nil {
		return "", err
	}

	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		// pkcs11.NewAttribute(pkcs11.CKA_ID, label),
	}
	handles, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return "", err
	}
	if len(handles) == 0 {
		return "", fmt.Errorf("NOT FIND KEY")
	}
	pvk := handles[0]
	fmt.Printf("Signing %d bytes with: %s\n", len(msg), msg)
	err = p.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil)}, pvk)
	if err != nil {
		return "", err
	}
	sig, err := p.Sign(session, []byte(msg))
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(sig)
	return signature, nil
}

func HSM_RSA_SHA256_Verify(hsmLib, pin, token, label string, data string, signature []byte) (bool, error) {
	// 1. Init HSM
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return false, err
	}

	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}

	pbks, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return false, err
	}
	if len(pbks) == 0 {
		return false, fmt.Errorf("NOT FOUND KEY")
	}
	pbk := pbks[0]
	err = p.VerifyInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil)}, pbk)
	if err != nil {
		return false, err
	}
	err = p.Verify(session, []byte(data), signature)
	if err != nil {
		return false, err
	}
	return true, nil
}

func HSM_RSA_Import_PrivateKey(hsmLib, pin, token, label string, privateKey *rsa.PrivateKey) error {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return err
	}
	id, err := hex.DecodeString(label)
	if err != nil {
		return err
	}

	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),

		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),

		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),

		pkcs11.NewAttribute(pkcs11.CKA_WRAP_WITH_TRUSTED, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),

		pkcs11.NewAttribute(pkcs11.CKA_MODULUS, privateKey.PublicKey.N.Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE_EXPONENT, big.NewInt(int64(privateKey.E)).Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_PRIME_1, new(big.Int).Set(privateKey.Primes[0]).Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_PRIME_2, new(big.Int).Set(privateKey.Primes[1]).Bytes()),

		pkcs11.NewAttribute(pkcs11.CKA_EXPONENT_1, new(big.Int).Set(privateKey.Precomputed.Dp).Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_EXPONENT_2, new(big.Int).Set(privateKey.Precomputed.Dq).Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_COEFFICIENT, new(big.Int).Set(privateKey.Precomputed.Qinv).Bytes()),
	}

	_, err = p.CreateObject(session, privateKeyTemplate)
	if err != nil {
		return err
	}

	return nil
}
func HSM_RSA_Import_PublicKey(hsmLib, pin, token, label string, publicKey *rsa.PublicKey) error {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return err
	}
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		// pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS, publicKey.N.Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, big.NewInt(int64(publicKey.E)).Bytes()),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, label),
	}
	_, err = p.CreateObject(session, publicKeyTemplate)
	if err != nil {
		return err
	}
	return nil
}

func HSM_RSA_Export_PublicKey(hsmLib, pin, token, label string) (*rsa.PublicKey, error) {
	p, cancel, session, err := HSM_Login(hsmLib, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}
	fmt.Printf("label: %s\n", label)
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, label),
	}
	pbks, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return nil, err
	}
	if len(pbks) == 0 {
		return nil, fmt.Errorf("NOT FOUND KEY")
	}
	pbk := pbks[0]
	pr, err := p.GetAttributeValue(session, pbk, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS, nil),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, nil),
	})
	if err != nil {
		return nil, err
	}

	modulus := new(big.Int)
	bigExponent := new(big.Int)
	exponent := int(bigExponent.SetBytes(pr[1].Value).Uint64())
	fmt.Printf("Export PublicKey: %d, %d \n", modulus, bigExponent)
	rsaPub := &rsa.PublicKey{
		N: modulus.SetBytes(pr[0].Value),
		E: exponent,
	}

	pubkeyPem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(rsaPub)}))
	fmt.Printf(" Public Key: \n%s\n", pubkeyPem)
	return rsaPub, nil
}
