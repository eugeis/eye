package core

import (
	"fmt"
	"errors"
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

func BuildAccessFinderFromFile(securityFile string) (ret AccessFinder, err error) {
	security := &Security{}
	ret = security
	err = fillAccessData(security, securityFile)
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
