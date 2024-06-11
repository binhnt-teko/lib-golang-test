package shamir

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/hashicorp/vault/shamir"
)

func GenKey(reader io.Reader) ([]byte, error) {
	buf := make([]byte, 2*aes.BlockSize)
	_, err := reader.Read(buf)
	return buf, err
}

func GenShares() ([][]byte, error) {
	SecretShares := int(5)
	SecretThreshold := int(3)

	newKey, err := GenKey(rand.Reader)
	if err != nil {
		fmt.Printf("GenKey failed: %s \n", err.Error())
		return nil, err
	}
	newKey64 := base64.StdEncoding.EncodeToString(newKey)
	fmt.Printf("rootKey:  %s \n", newKey64)

	shares, err := shamir.Split(newKey, SecretShares, SecretThreshold)
	if err != nil {
		fmt.Printf("shamir.Split failed: %s \n", err.Error())
		return nil, err
	}
	return shares, err
}
func VerifyShares(RekeyProgress [][]byte) {
	// var RekeyProgress [][]byte
	// for _, key := range shares {
	// 	RekeyProgress = append(RekeyProgress, []byte(key))
	// }

	recoveredKey, err := shamir.Combine(RekeyProgress)
	if err != nil {
		fmt.Printf("Error in Combine: %s \n", err.Error())
		return
	}
	newKey64 := base64.StdEncoding.EncodeToString(recoveredKey)
	fmt.Printf("recoveredKey:  %s \n", newKey64)
}
