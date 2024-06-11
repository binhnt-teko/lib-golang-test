package util

import (
	"time"

	"github.com/hashicorp/vault-client-go"
)

func GetValidClient(vaultUrl string, token string) (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress(vaultUrl),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(token); err != nil {
		return nil, err
	}
	return client, nil
}
