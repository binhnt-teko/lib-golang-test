package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

func Vault_KV_Get(client *vault.Client, ctx context.Context, mount_path string, path string) (map[string]interface{}, error) {
	// read the secret
	s, err := client.Secrets.KvV2Read(ctx, path, vault.WithMountPath(mount_path))
	if err != nil {
		return nil, err
	}
	return s.Data.Data, nil
}

func Vault_KV_Set(client *vault.Client, ctx context.Context, mount_path string, path string, value map[string]any) error {
	// write a secret
	_, err := client.Secrets.KvV2Write(ctx, path, schema.KvV2WriteRequest{
		Data: value},
		vault.WithMountPath(mount_path),
	)
	if err != nil {
		fmt.Printf("KvV2Write failed %s \n", err.Error())
		return err
	}
	fmt.Println("secret written successfully")
	return nil
}
