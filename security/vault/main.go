package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blcvn/lib-golang-test/security/vault/util"
)

func main() {
	vaultUrl := "https://vault.epr.vn"
	vaultToken := os.Getenv("TOKEN")
	client, err := util.GetValidClient(vaultUrl, vaultToken)
	if err != nil {
		fmt.Printf("GetValidClient failed: %s \n", err.Error())
		return
	}

	mountPath := "test"
	ctx := context.Background()

	path_test := "test"
	data_test := map[string]any{
		"password": "12323123",
		"sender":   "abc",
		"receiver": "def",
	}
	err = util.Vault_KV_Set(client, ctx, mountPath, path_test, data_test)
	if err != nil {
		fmt.Printf("Vault_KV_Set failed %s \n", err.Error())
		return
	}
	fmt.Printf("Write successfully")
}
