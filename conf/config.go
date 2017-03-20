package conf

import (
	"github.com/jinzhu/configor"
	"rest/service"
	"fmt"
	"rest/integ"
)

var log = integ.Log

type RestConfig struct {
	Name  string `default:"app name"`
	Port  int `default:"3000"`
	Debug bool `default:true`

	MySql []*service.MySql
	Http  []*service.Http

	serviceFactory service.Factory
}

func Load(fileNames string) *RestConfig {
	ret := &RestConfig{}
	configor.Load(ret, fileNames)
	return ret
}

func (c *RestConfig) Print() {
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

func (c *RestConfig) ServiceFactory() service.Factory {
	if c.serviceFactory == nil {
		serviceFactory := service.NewFactory()
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
