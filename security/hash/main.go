package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/blcvn/lib-golang-test/security/hash/util"
)

const MAX_TEST = uint64(10000000)

func Test_Hash() {
	h := sha256.New()

	start := time.Now()
	for i := range MAX_TEST {
		var s = big.NewInt(int64(i)) // int to big Int

		h.Write(s.Bytes())
		h.Sum(nil)
		// fmt.Printf("%x\n", data)
	}
	fmt.Printf("-- Hash: %d ms\n", time.Since(start).Milliseconds())
}
func Test_AES_Encrypt() {
	password := "THU NGHIEM"
	h := md5.New()
	h.Write([]byte(password))
	key := h.Sum(nil)

	start := time.Now()
	for i := range MAX_TEST {
		var s = big.NewInt(int64(i)) // int to big Int

		util.EncryptMessage(key, s.Bytes())
		// fmt.Printf("%x\n", data)
	}
	fmt.Printf("-- AES_Encrypt: %d ms\n", time.Since(start).Milliseconds())
}
func main() {
	Test_Hash()
	Test_AES_Encrypt()
}
