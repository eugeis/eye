package core

import (
	"strings"
	"fmt"
	"errors"
	"path/filepath"
)

type Access struct {
	User     string
	Password []byte
}

type AccessFinder interface {
	FindAccess(key string) (Access, error)
}

type Security struct {
	Access map[string]Access
}

func (o *Security) FindAccess(key string) (ret Access, err error) {
	var ok bool
	ret, ok = o.Access[key]
	if !ok {
		err = errors.New(fmt.Sprintf("No access data found for '%v'", key))
	}
	return
}

func BuildAccessFinder(config *Config) (ret AccessFinder, err error) {
	security := &Security{}
	ret = security
	if strings.EqualFold(config.SecurityType.Type, "file") {
		path := filepath.Dir(config.ConfigFile)
		err = fillAccessData(security, filepath.Join(path, config.SecurityType.File))
	} else if strings.EqualFold(config.SecurityType.Type, "console") {
		security.Access = extractAccessKeys(config)
		fillAccessDataFromConsole(security)

	} else if strings.EqualFold(config.SecurityType.Type, "vault") {
		security.Access = extractAccessKeys(config)
		var vault *VaultClient
		vault, err = NewVaultClient(config.SecurityType.Token)
		if err == nil {
			vault.fillAccessData(config.Name, security)
		}
	}
	return
}

func extractAccessKeys(config *Config) (ret map[string]Access) {
	ret = make(map[string]Access)
	for _, item := range config.MySql {
		if _, ok := ret[item.AccessKey]; !ok {
			ret[item.AccessKey] = Access{}
		}
	}
	for _, item := range config.Http {
		if _, ok := ret[item.AccessKey]; !ok {
			ret[item.AccessKey] = Access{}
		}
	}
	return
}
