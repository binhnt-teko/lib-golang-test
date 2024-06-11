package main

import (
	"encoding/base64"
	"fmt"

	"github.com/blcvn/lib-golang-test/zk/pair/shamir"
)

func main() {
	// rand.Seed(time.Now().UnixNano())

	// fmt.Println(util.IsPrime(11, 5)) // true
	// fmt.Println(util.IsPrime(20, 5)) // false

	// util.Ecliptic_Curve_Arithmetic()
	shares, err := shamir.GenShares()
	if err != nil {
		return
	}
	for index, share := range shares {
		shareKey := base64.StdEncoding.EncodeToString(share)
		fmt.Printf("Key[%d]:  %s \n", index, shareKey)
	}
	// fmt.Println("%s", string(recoveredKey))

	shamir.VerifyShares(shares)
}
