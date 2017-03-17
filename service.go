package main

import (
	"regexp"
	"strings"
)

type Service interface {
	Name() string
	Kind() string

	Init() error
	Close()

	Ping() error
	Check(query string, pattern *regexp.Regexp) (bool, error)
}

type ServiceFactory interface {
	FindCheck(name string) *Service
	Close()
}

type ServiceFactoryByConfig struct {
	Config   *Config
	services map[string]Service
}

func (sf *ServiceFactoryByConfig) FindService(name string) Service {
	if sf.services == nil {
		sf.services = make(map[string]Service)
	}
	nameLowCase := strings.ToLower(name)
	var s Service = sf.services[nameLowCase]
	if s == nil {
		for _, v := range sf.Config.MySql {
			if strings.EqualFold(v.Name(), name) {
				if v.Init() == nil {
					s = v
					sf.services[nameLowCase] = s
				}
				break
			}
		}
	}

	if s == nil {
		for _, v := range sf.Config.Http {
			if strings.EqualFold(v.Name(), name) {
				if v.Init() == nil {
					s = v
					sf.services[nameLowCase] = s
				}
				break
			}
		}
	}
	return s
}

func (sf *ServiceFactoryByConfig) Close() {
	for _, v := range sf.services {
		v.Close()
	}
}
