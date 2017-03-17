package main

import (
	"github.com/jinzhu/configor"
	"rest/service"
	"fmt"
)

type Config struct {
	Name  string `default:"app name"`
	Port  int `default:"3000"`
	Debug bool `default:true`

	MySql []*service.MySql
	Http  []*service.Http
}

func load(fileNames string) *Config {
	ret := &Config{}
	configor.Load(ret, fileNames)
	return ret
}

func (c *Config) print() {
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
