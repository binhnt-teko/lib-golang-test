package util

import (
	"fmt"

	"github.com/miekg/pkcs11"
)

func HSM_ECDSA_SHA256_Sign(libPath string, pin string, token string, signatureInput []byte) ([]byte, error) {
	//1. Init HSM
	p, cancel, session, err := HSM_Login(libPath, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
	}
	handles, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return nil, err
	}
	privateKeyHandle := handles[0]

	// Digest
	err = p.DigestInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256, nil)})
	if err != nil {
		return nil, fmt.Errorf("DigestInit failed: %s", err)
	}
	hash, err := p.Digest(session, signatureInput)
	if err != nil {
		return nil, fmt.Errorf("Digest failed: %s", err)
	}

	//Signature
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}
	err = p.SignInit(session, mechanism, privateKeyHandle)
	if err != nil {
		return nil, fmt.Errorf("pkcs11key: sign init: %s", err)
	}
	signature, err := p.Sign(session, hash)
	if err != nil {
		return nil, fmt.Errorf("pkcs11key: sign: %s", err)
	}
	signatureRFC, err := ecdsaPKCS11ToRFC5480(signature)
	if err != nil {
		return nil, fmt.Errorf("ecdsaPKCS11ToRFC5480 failed: %s", err)
	}
	return signatureRFC, nil
}
