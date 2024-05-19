package main

import (
	usecase "github.com/blcvn/lib-golang-test/security/hsm/usecase"
)

func main() {
	usecase.HSM_Exchange_AES_Key()
	// usecases.AES_CFB()
	// usecases.LoadConfig_FileKey()
	// usecases.LoadConfig_HSM()
	// usecases.GenSKI()
	// usecases.ReadObject()
	// usecases.TestImportPrivateKey()
	// usecases.TestSignAndVerify()
	// usecases.TestAES_CBC()
}
