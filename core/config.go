package core

import (
	"github.com/jinzhu/configor"
	"fmt"
	"rest/integ"
)

var log = integ.Log

type Config struct {
	Name  string `default:"app Name"`
	Port  int `default:"3000"`
	Debug bool `default:true`
	Token string

	MySql []*MySql
	Http  []*Http

	serviceFactory Factory
}

func Load(fileNames string) *Config {
	ret := &Config{}
	configor.Load(ret, fileNames)
	return ret
}

func (c *Config) Print() {
	log.Info(fmt.Sprintf("Name: %s, Port: %d, Debug: %v", c.Name, c.Port, c.Debug))
	log.Info("MySql:")
	for _, v := range c.MySql {
		log.Info(fmt.Sprintf("\t%s", v.Name()))
	}

	log.Info("Http:")
	for _, v := range c.Http {
		log.Info(fmt.Sprintf("\t%s", v.Name()))
	}
}

func (c *Config) ServiceFactory() Factory {
	if c.serviceFactory == nil {
		serviceFactory := NewFactory()
		for _, service := range c.MySql {
			serviceFactory.Add(service)
		}

		for _, service := range c.Http {
			serviceFactory.Add(service)
		}
		c.serviceFactory = &serviceFactory
	}
	return c.serviceFactory
}
