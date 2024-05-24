package main

import (
	"fmt"

	"github.com/thedonutfactory/go-tfhe/gates"
)

func main() {
	// generate public and private keys
	ctx := gates.DefaultGateBootstrappingParameters(100)
	pub, prv := ctx.GenerateKeys()

	// encrypt 2 8-bit ciphertexts
	x := prv.Encrypt(int8(22))
	y := prv.Encrypt(int8(33))

	// perform homomorphic sum gate operations
	BITS := 8
	temp := ctx.Int(3)
	sum := ctx.Int(9)
	carry := ctx.Int2()
	for i := 0; i < BITS; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = pub.Xor(x[i], y[i]) // temp = xi XOR yi
		sum[i] = pub.Xor(temp[0], carry[0])

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		temp[1] = pub.And(x[i], y[i])
		temp[2] = pub.And(carry[0], temp[0])
		carry[1] = pub.Xor(temp[1], temp[2])
		carry[0] = pub.Copy(carry[1])
	}
	sum[BITS] = pub.Copy(carry[0])

	// decrypt results
	z := prv.Decrypt(sum[:])
	fmt.Println("The sum of of x and y: ", z)
}
