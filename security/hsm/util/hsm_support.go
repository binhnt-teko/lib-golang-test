package util

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
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

func ecdsaPKCS11ToRFC5480(pkcs11Signature []byte) (rfc5480Signature []byte, err error) {
	mid := len(pkcs11Signature) / 2

	r := &big.Int{}
	s := &big.Int{}

	return asn1.Marshal(rfc5480ECDSASignature{
		R: r.SetBytes(pkcs11Signature[:mid]),
		S: s.SetBytes(pkcs11Signature[mid:]),
	})
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

func verifyANS1(publicKey *ecdsa.PublicKey, signature []byte, hash []byte) bool {
	var esig struct {
		R, S *big.Int
	}
	if _, err := asn1.Unmarshal(signature, &esig); err != nil {
		fmt.Println("asn1.Unmarshal error: ", err)
		return false
	}
	return ecdsa.Verify(publicKey, hash, esig.R, esig.S)

}
