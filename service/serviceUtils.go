package service

import (
	"regexp"
	"time"
	"context"
	"rest/integ"
	"strings"
)

var log = integ.Log

func Match(info string, pattern *regexp.Regexp, query func() ([]byte, error)) (ok bool, err error) {
	data, err := query()
	if err != nil {
		log.Info("Service failed because of '%s'", err)
	} else {
		log.Debug("data '%s'", data)
		ok = pattern.Match(data)
		if err != nil {
			log.Info("error '%s' for %s", info)
		} else if !ok {
			log.Debug("no match for %s", info)
		} else {
			log.Debug("match for %s", info)
		}
	}
	return
}

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

func (o *SimpleServiceFactory) Find(name string) Service {
	nameLowCase := strings.ToLower(name)
	return o.services[nameLowCase]
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
