package core

import vaultapi "github.com/hashicorp/vault/api"

type VaultClient struct {
	client  *vaultapi.Client
	logical *vaultapi.Logical
}

func NewVaultClient(token string) (ret *VaultClient, err error) {
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err == nil {
		client.SetToken(token)
		ret = &VaultClient{client: client, logical: client.Logical()}
	}
	return
}

func (o *VaultClient) Access(key string) (ret Access, err error) {
	var secret *vaultapi.Secret
	secret, err = o.logical.Read("eye/" + key)
	if err == nil {
		ret = Access{Key: key, User: secret.Data["user"].(string), Password: secret.Data["password"].(string)}
	}
	return
}
