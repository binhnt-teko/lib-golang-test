package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/miekg/pkcs11"
)

type rfc5480ECDSASignature struct {
	R, S *big.Int
}

var (
	// OAEPLabel defines the label we use for OAEP encryption; this cannot be changed
	OAEPLabel = []byte("")

	// OAEPSha1Params describes the OAEP parameters with sha1 hash algorithm; needed by SoftHSM
	OAEPSha1Params = &pkcs11.OAEPParams{
		HashAlg:    pkcs11.CKM_SHA_1,
		MGF:        pkcs11.CKG_MGF1_SHA1,
		SourceType: pkcs11.CKZ_DATA_SPECIFIED,
		SourceData: OAEPLabel,
	}
	// OAEPSha256Params describes the OAEP parameters with sha256 hash algorithm
	OAEPSha256Params = &pkcs11.OAEPParams{
		HashAlg:    pkcs11.CKM_SHA256,
		MGF:        pkcs11.CKG_MGF1_SHA256,
		SourceType: pkcs11.CKZ_DATA_SPECIFIED,
		SourceData: OAEPLabel,
	}
)

func HSM_ListSlot(libPath string) error {
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		fmt.Printf("Initialize failed %s \n", err.Error())
		return err
	}
	defer p.Destroy()
	defer p.Finalize()
	slots, err := p.GetSlotList(true)
	if err != nil {
		fmt.Printf("GetSlotList failed %s \n", err.Error())
		return err
	}
	for _, slot := range slots {
		slotInfo, err := p.GetSlotInfo(slot)
		if err != nil {
			continue
		}
		fmt.Printf("Slot: %s\n", slotInfo.SlotDescription)
	}
	return nil
}

func HSM_Decrypt(libPath string, slotID uint, pin string, numKey int, sDec []byte) ([]byte, error) {
	//1. Init HSM
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		fmt.Printf("Initialize failed %s \n", err.Error())
		return nil, err
	}
	defer p.Destroy()
	defer p.Finalize()

	session, err := p.OpenSession(slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		fmt.Printf("OpenSession failed %s \n", err.Error())
		return nil, err
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		fmt.Printf("Login failed %s \n", err.Error())
		return nil, err
	}
	defer p.Logout(session)

	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
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
	if err := p.FindObjectsInit(session, findTemplate); err != nil {
		p.Destroy()
		fmt.Printf("FindObjectsInit failed: %s \n", err)
		return nil, err
	}
	handles, moreAvailable, err := p.FindObjects(session, numKey)
	if err != nil {
		p.Destroy()
		fmt.Printf("FindObjects failed: %s \n", err)
		return nil, err
	}
	if moreAvailable {
		fmt.Printf("Too many object return from token \n")
		return nil, fmt.Errorf("NUMBER OBJECT OVER %d", numKey)
	}
	if len(handles) == 0 {
		fmt.Printf("Cannot find hanndles \n")
		return nil, fmt.Errorf("NOT FIND PRIVATE KEY")
	}

	if err = p.FindObjectsFinal(session); err != nil {
		p.Destroy()
		fmt.Printf("FindObjectsFinal failed: %s \n", err.Error())
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
func getSlotId(p *pkcs11.Ctx, slots []uint, token string) uint {
	slotID := slots[0]
	for _, slot := range slots {
		tokenInfo, err := p.GetTokenInfo(slot)
		if err != nil {
			continue
		}
		if tokenInfo.Label == token {
			slotID = slot
			break
		}
	}
	return slotID
}
func HSM_ImportKey(libPath string, pin string, token string, privateKey *ecdsa.PrivateKey) error {
	//1. Init
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		fmt.Printf("HSM_ImportKey: Initialize failed %s \n", err.Error())
		return err
	}
	defer p.Destroy()
	defer p.Finalize()

	//2 . Find token
	slots, err := p.GetSlotList(true)
	if err != nil {
		fmt.Printf("HSM_ImportKey: GetSlotList failed %s \n", err.Error())
		return err
	}
	if len(slots) == 0 {
		fmt.Printf("HSM_ImportKey: GetSlotList not find slots")
		return fmt.Errorf("NOT FIND SLOTS")
	}
	slotID := getSlotId(p, slots, token)
	fmt.Printf("HSM_ImportKey: slotID %d\n", slotID)

	session, err := p.OpenSession(slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		fmt.Printf("HSM_ImportKey: OpenSession failed %s \n", err.Error())
		return err
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		fmt.Printf("HSM_ImportKey: Login failed %s \n", err.Error())
		return err
	}
	defer p.Logout(session)

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

func HSM_Sign(libPath string, pin string, token string, signatureInput []byte) ([]byte, error) {
	//1. Init
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		fmt.Printf("Initialize failed %s \n", err.Error())
		return nil, err
	}
	defer p.Destroy()
	defer p.Finalize()

	//2 . Find token
	slots, err := p.GetSlotList(true)
	if err != nil {
		fmt.Printf("GetSlotList failed %s \n", err.Error())
		return nil, err
	}
	if len(slots) == 0 {
		fmt.Printf("GetSlotList not find slots")
		return nil, fmt.Errorf("NOT FIND SLOTS")
	}
	slotID := getSlotId(p, slots, token)
	fmt.Printf("Slot: %d\n", slotID)
	session, err := p.OpenSession(slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		fmt.Printf("OpenSession failed %s \n", err.Error())
		return nil, err
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		fmt.Printf("Login failed %s \n", err.Error())
		return nil, err
	}
	defer p.Logout(session)

	//Find private key
	searchTemplates := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		// pkcs11.NewAttribute(pkcs11.CKA_ID, publicKeyID),
	}
	if err := p.FindObjectsInit(session, searchTemplates); err != nil {
		return nil, err
	}

	handles, moreAvailable, err := p.FindObjects(session, 1)
	if err != nil {
		return nil, err
	}
	if moreAvailable {
		return nil, errors.New("too many objects returned from FindObjects")
	}
	if err = p.FindObjectsFinal(session); err != nil {
		return nil, err
	} else if len(handles) == 0 {
		return nil, errors.New("no objects found")
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

func ecdsaPKCS11ToRFC5480(pkcs11Signature []byte) (rfc5480Signature []byte, err error) {
	mid := len(pkcs11Signature) / 2

	r := &big.Int{}
	s := &big.Int{}

	return asn1.Marshal(rfc5480ECDSASignature{
		R: r.SetBytes(pkcs11Signature[:mid]),
		S: s.SetBytes(pkcs11Signature[mid:]),
	})
}

func HSM_ReadObject(libPath string, token string, pin string) ([]byte, error) {
	//1. Init
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		fmt.Printf("Initialize failed %s \n", err.Error())
		return nil, err
	}
	defer p.Destroy()
	defer p.Finalize()

	//2 . Find token
	slots, err := p.GetSlotList(true)
	if err != nil {
		fmt.Printf("GetSlotList failed %s \n", err.Error())
		return nil, err
	}
	if len(slots) == 0 {
		fmt.Printf("GetSlotList not find slots")
		return nil, fmt.Errorf("NOT FIND SLOTS")
	}
	slotID := getSlotId(p, slots, token)
	session, err := p.OpenSession(slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		fmt.Printf("OpenSession failed %s \n", err.Error())
		return nil, err
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		fmt.Printf("Login failed %s \n", err.Error())
		return nil, err
	}
	defer p.Logout(session)

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
	if err := p.FindObjectsInit(session, findTemplate); err != nil {
		p.Destroy()
		fmt.Printf("FindObjectsInit failed: %s \n", err)
		return nil, err
	}
	numKey := 10
	handles, moreAvailable, err := p.FindObjects(session, numKey)
	if err != nil {
		p.Destroy()
		fmt.Printf("FindObjects failed: %s \n", err)
		return nil, err
	}
	if moreAvailable {
		fmt.Printf("Too many object return from token \n")
		return nil, fmt.Errorf("NUMBER OBJECT OVER %d", numKey)
	}
	if len(handles) == 0 {
		fmt.Printf("Cannot find hanndles \n")
		return nil, fmt.Errorf("NOT FIND PRIVATE KEY")
	}

	if err = p.FindObjectsFinal(session); err != nil {
		p.Destroy()
		fmt.Printf("FindObjectsFinal failed: %s \n", err.Error())
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

func HSM_VerifyANS1(publicKey *ecdsa.PublicKey, signature []byte, hash []byte) bool {
	var esig struct {
		R, S *big.Int
	}
	if _, err := asn1.Unmarshal(signature, &esig); err != nil {
		fmt.Println("asn1.Unmarshal error: ", err)
		return false
	}
	return ecdsa.Verify(publicKey, hash, esig.R, esig.S)

}
func getId(num uint16) ([]byte, error) {
	// first lookup the key
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	id := buf.Bytes()
	return id, nil
}

func HSM_AES_Key(libPath string, pin string, token string) error {
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		return err
	}

	defer p.Destroy()
	defer p.Finalize()

	slots, err := p.GetSlotList(true)
	if err != nil {
		return err
	}
	slotId := getSlotId(p, slots, token)
	session, err := p.OpenSession(slotId, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return err
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		return err
	}
	defer p.Logout(session)

	id := hex.EncodeToString([]byte(token))

	aesKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, false),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false), // we don't need to extract this..
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, make([]byte, 32)), /* KeyLength */
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, token),            /* Name of Key */
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),
	}

	if err := p.FindObjectsInit(session, aesKeyTemplate); err != nil {
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

	aesKey, err := p.CreateObject(session, aesKeyTemplate)
	if err != nil {
		return err
	}
	log.Printf("Created AES Key: %v", aesKey)
	return nil
}

func HSM_AES_CBC_Encrypt(libPath string, pin string, token string, message []byte) ([]byte, error) {
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		return nil, err
	}

	defer p.Destroy()
	defer p.Finalize()

	slots, err := p.GetSlotList(true)
	if err != nil {
		return nil, err
	}
	slotId := getSlotId(p, slots, token)
	session, err := p.OpenSession(slotId, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, err
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		return nil, err
	}
	defer p.Logout(session)

	id := hex.EncodeToString([]byte(token))

	ktemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, token),
	}
	if err := p.FindObjectsInit(session, ktemplate); err != nil {
		return nil, err
	}
	secretHandles, _, err := p.FindObjects(session, 100)
	if err != nil {
		return nil, err
	}
	if err = p.FindObjectsFinal(session); err != nil {
		return nil, err
	}
	if len(secretHandles) == 0 {
		return nil, fmt.Errorf("NOT FOUND KEYS")
	}

	iv := make([]byte, 16)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Handle.Len: %d\n", len(secretHandles))
	fmt.Printf("Handle: %d\n", secretHandles[0])
	err = p.EncryptInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)}, secretHandles[0])
	if err != nil {
		return nil, err
	}

	ct, err := p.Encrypt(session, message)
	if err != nil {
		fmt.Printf("Encrypt() failed %s\n", err)
		return nil, err
	}
	cdWithIV := append(iv, ct...)

	fmt.Println("Test encrypt............")
	iv1 := cdWithIV[0:16]
	fmt.Printf("iv1: %s\n", hex.EncodeToString(iv1))

	err = p.DecryptInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)}, secretHandles[0])
	if err != nil {
		fmt.Printf("DecryptInit() failed %s\n", err)
		return nil, err
	}
	clear_text, err := p.Decrypt(session, cdWithIV[16:])
	if err != nil {
		fmt.Printf("Decrypt() failed %s\n", err)
		return nil, err
	}
	fmt.Println(string(clear_text))

	// append the IV to the ciphertext

	fmt.Printf("iv: %s\n", hex.EncodeToString(iv))
	//base64.RawStdEncoding.EncodeToString(cdWithIV)
	return cdWithIV, nil

}

func HSM_AES_CBC_Decrypt(libPath string, pin string, token string, encrypted []byte) ([]byte, error) {
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		return nil, err
	}

	defer p.Destroy()
	defer p.Finalize()

	slots, err := p.GetSlotList(true)
	if err != nil {
		return nil, err
	}
	slotId := getSlotId(p, slots, token)
	session, err := p.OpenSession(slotId, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		panic(err)
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		panic(err)
	}
	defer p.Logout(session)

	id := hex.EncodeToString([]byte(token))

	ktemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, token),
	}
	if err := p.FindObjectsInit(session, ktemplate); err != nil {
		return nil, err
	}
	secretHandles, _, err := p.FindObjects(session, 100)
	if err != nil {
		return nil, err
	}
	if err = p.FindObjectsFinal(session); err != nil {
		return nil, err
	}
	if len(secretHandles) == 0 {
		return nil, fmt.Errorf("NOT FOUND KEYS")
	}
	fmt.Printf("Handle.Len: %d\n", len(secretHandles))
	fmt.Printf("Handle: %d\n", secretHandles[0])

	iv := encrypted[0:16]
	fmt.Printf("iv: %s\n", hex.EncodeToString(iv))

	err = p.DecryptInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)}, secretHandles[0])
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
