package util

import (
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	"github.com/miekg/pkcs11"
)

func HSM_ImportKey(libPath string, pin string, token string, privateKey *ecdsa.PrivateKey) error {
	//1. Init HSM
	p, cancel, session, err := HSM_Login(libPath, pin, token)
	defer cancel(session)
	if err != nil {
		return err
	}

	oidNamedCurveP256 := asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	params, err := asn1.Marshal(oidNamedCurveP256)
	if err != nil {
		fmt.Printf("asn1.Marshal sha256 failed %s \n", err.Error())
		return err
	}
	value, err := asn1.Marshal(privateKey.D)
	if err != nil {
		fmt.Printf("asn1.Marshal D failed %s \n", err.Error())
		return err
	}
	tokenId := hex.EncodeToString([]byte(token))

	// point, err := asn1.Marshal(privateKey.PublicKey.)
	// if err != nil {
	// 	fmt.Printf("asn1.Marshal D failed %s \n", err.Error())
	// 	return err
	// }

	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, token), //label
		// pkcs11.NewAttribute(pkcs11.CKA_SUBJECT, tokenId),
		pkcs11.NewAttribute(pkcs11.CKA_ID, tokenId),
		// pkcs11.NewAttribute(pkcs11.CKA_DERIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		// pkcs11.NewAttribute(pkcs11.CKA_VALUE, make([]byte, 32)), /* KeyLength */
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY_RECOVER, true),
		// pkcs11.NewAttribute(pkcs11.CKA_WRAP_WITH_TRUSTED, true),
		// pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		// pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, params),
		// pkcs11.NewAttribute(pkcs11.CKA_ECDSA_PARAMS, []byte{}),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, value),
		// pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, point),
	}

	if err := p.FindObjectsInit(session, privateKeyTemplate); err != nil {
		return err
	}
	secretHandles, _, err := p.FindObjects(session, 100)

	if err != nil {
		return err
	}
	if err = p.FindObjectsFinal(session); err != nil {
		return err
	}
	if len(secretHandles) > 0 {
		return fmt.Errorf("KEY FOUND")
	}

	_, err = p.CreateObject(session, privateKeyTemplate)
	if err != nil {
		fmt.Printf("CreateObject failed: %s \n", err.Error())
		return err
	}

	return nil
}

func HSM_ReadObject(libPath string, token string, pin string) ([]byte, error) {
	//1. Init
	p, cancel, session, err := HSM_Login(libPath, pin, token)
	defer cancel(session)
	if err != nil {
		return nil, err
	}

	findTemplate := []*pkcs11.Attribute{
		// pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		// pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		// pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		// pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		// pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		// pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		// pkcs11.NewAttribute(pkcs11.CKA_WRAP, false),
		// pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, false),
		// pkcs11.NewAttribute(pkcs11.CKA_VERIFY, false),
		// pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		// pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true), // we do need to extract this
		// pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		// pkcs11.NewAttribute(pkcs11.CKA_VALUE, make([]byte, 32)), /* KeyLength */
		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, "AESKeyToWrap"),   /* Name of Key */
		// pkcs11.NewAttribute(pkcs11.CKA_ID, id),
	}
	handles, err := HSM_Find(p, session, findTemplate, 1)
	if err != nil {
		return nil, err
	}

	searchTemplates := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, "priv1"),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SUBJECT, []byte{}),
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte{}),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, false),
		pkcs11.NewAttribute(pkcs11.CKA_DERIVE, true),

		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),

		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),

		pkcs11.NewAttribute(pkcs11.CKA_WRAP_WITH_TRUSTED, false),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, false),

		// pkcs11.NewAttribute(pkcs11.CKA_GOSTR3410_PARAMS, privateKey),
		// Parameters ::= CHOICE {
		// 	ecParameters  ECParameters,
		// 	namedCurve    CURVES.&id({CurveNames}),
		// 	implicitlyCA  NULL
		//   }
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, []byte{}),
		pkcs11.NewAttribute(pkcs11.CKA_ECDSA_PARAMS, []byte{}),
		// pkcs11.NewAttribute(pkcs11.CKA_VALUE, []byte{}),
	}

	for _, hd := range handles {
		atrs, err := p.GetAttributeValue(session, hd, searchTemplates)
		if err != nil {
			fmt.Printf("GetAttributeValue failed %s \n", err.Error())
			return nil, err
		}
		for _, attr := range atrs {
			fmt.Printf("%08x: %+v\n", attr.Type, attr.Value)
			if attr.Type == pkcs11.CKA_LABEL {
				fmt.Printf("CKA_LABEL: %s\n", string(attr.Value))
			}
			if attr.Type == pkcs11.CKA_ID {
				fmt.Printf("CKA_ID: %s\n", hex.EncodeToString(attr.Value))
			}
			if attr.Type == pkcs11.CKA_EC_PARAMS {
				fmt.Printf("%s\n", hex.EncodeToString(attr.Value))
				t := asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
				t1, _ := asn1.Marshal(t)
				fmt.Printf("T: %s\n", hex.EncodeToString(t1))

			}
		}
	}
	return nil, nil
}
