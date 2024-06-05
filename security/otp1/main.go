package main

import (
	"fmt"

	"github.com/xlzd/gotp"
)

func main() {
	totp := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO")
	otpCode := totp.Now() // current otp
	fmt.Printf("Current otp: %s \n", otpCode)
	// otp of timestamp 1524486261 '123456'

	otpCode1 := totp.At(1524486261)
	fmt.Printf("otp at  1524486261: %s \n", otpCode1)

	//  OTP verified for a given timestamp
	if totp.Verify(otpCode1, 1524486261) {
		fmt.Printf("Verify at  1524486261: %s \n", otpCode1)
	}
	if !totp.Verify(otpCode1, 1524486800) {
		fmt.Printf("Verify at  1524486800: %s => failed \n", otpCode1)
	}

	// generate a provisioning uri
	otpauth := totp.ProvisioningUri("demoAccountName", "issuerName")
	fmt.Printf("otpauth: %s \n", otpauth)
	// otpauth://totp/issuerName:demoAccountName?secret=4S62BZNFXXSZLCRO&issuer=issuerName
}
