package core

import (
	"time"
	"context"
	"strings"
	"errors"
	"fmt"
)

func TimeoutContext(timeout time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), timeout)
	return c
}

type SimpleServiceFactory struct {
	services map[string]Service
}

func NewFactory() SimpleServiceFactory {
	return SimpleServiceFactory{services: make(map[string]Service)}
}

func (o *SimpleServiceFactory) Find(name string) (ret Service, err error) {
	nameLowCase := strings.ToLower(name)
	ret = o.services[nameLowCase]
	if ret == nil {
		err = errors.New(fmt.Sprintf("Can't find the service '%s'", name))
	}
	return
}

func (o *SimpleServiceFactory) Add(service Service) {
	nameLowCase := strings.ToLower(service.Name())
	o.services[nameLowCase] = service
}

func (o *SimpleServiceFactory) Close() {
	for _, item := range o.services {
		item.Close()
	}
	o.services = make(map[string]Service)
}

type Access struct {
	Key string
	User string
	Password string
}