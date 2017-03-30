package core

import (
	vaultapi "github.com/hashicorp/vault/api"
	"fmt"
	"strings"
)

type VaultClient struct {
	client  *vaultapi.Client
	logical *vaultapi.Logical
}

func NewVaultClient(token string, address string) (ret *VaultClient, err error) {
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err == nil {
		client.SetToken(token)
		if len(address) > 0 {
			client.SetAddress(address)
		}
		ret = &VaultClient{client: client, logical: client.Logical()}
	}
	return
}

func (o *VaultClient) fillAccessData(name string, security *Security) (err error) {
	var secret *vaultapi.Secret
	basePath := fmt.Sprintf("secret/eye/%v", strings.ToLower(name))
	for key, item := range security.Access {
		vaultPath := fmt.Sprintf("%v/%v", basePath, key)
		secret, err = o.logical.Read(vaultPath)
		if err != nil {
			break
		}
		if secret != nil {
			item.User = secret.Data["user"].(string)
			item.Password = []byte(secret.Data["password"].(string))
			security.Access[key] = item
		} else {
			l.Info("No access data in Vault for '%key", vaultPath)
		}

	}
	return
}

func BuildAccessFinderFromVault(vaultToken string, vaultAddress string, config *Config) (ret AccessFinder, err error) {

	security := &Security{}
	ret = security

	security.Access = extractAccessKeys(config)
	var vault *VaultClient
	vault, err = NewVaultClient(vaultToken, vaultAddress)
	if err == nil {
		err = vault.fillAccessData(config.Name, security)
	}
	return
}
