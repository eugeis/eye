package core

import (
	vaultapi "github.com/hashicorp/vault/api"
	"fmt"
)

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

func (o *VaultClient) fillAccessData(name string, security *Security) {
	for key, item := range security.Access {
		secret, err := o.logical.Read(fmt.Sprintf("eye/%v/%v", name, key))
		if err == nil {
			item.User = secret.Data["user"].(string)
			item.Password = []byte(secret.Data["password"].(string))
			security.Access[key] = item
		}
	}
}
