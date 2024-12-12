package main

import (
	"github.com/blcvn/lib-golang-test/mpc/tss/dwjpeng"
	// "github.com/blcvn/lib-golang-test/mpc/tss/bnb"
	// multipartysig "github.com/blcvn/lib-golang-test/mpc/tss/multi-party-sig"
	// "github.com/blcvn/lib-golang-test/mpc/tss/frost-ed25519"
	// ibm "github.com/blcvn/lib-golang-test/mpc/tss/tss-ibm"
)

func main() {
	dwjpeng.RUN_TSS_ECDSA()
	// bnb.RUN_TSS_ECDSA()
	// multipartysig.RUN_TSS_ECDSA()
	// frost.RUN_TSS()
	// ibm.RUN_TSS()
}
