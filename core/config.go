package core

import (
	"github.com/jinzhu/configor"
	"fmt"
	"eye/integ"
	"strings"
)

var log = integ.Log

type Security struct {
	Access []Access
}

type Config struct {
	Name  string `default:"app Name"`
	Port  int `default:"3000"`
	Debug bool `default:true`
	Token string

	MySql []*MySql
	Http  []*Http

	fileNames []string
}

func Load(fileNames []string) (ret *Config, err error) {
	ret = &Config{fileNames: fileNames}
	err = configor.Load(ret, fileNames...)

	//ignore, https://github.com/jinzhu/configor/issues/6
	if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
		err = nil
	}
	return
}

func LoadSecurity(fileNames []string) (ret *Security, err error) {
	ret = &Security{}
	err = configor.Load(ret, fileNames...)

	//ignore, https://github.com/jinzhu/configor/issues/6
	if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
		err = nil
	}
	return
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

func (c *Security) FindAccess(key string) (ret Access) {
	for _, access := range c.Access {
		if strings.EqualFold(key, access.Key) {
			ret = access
			break
		}
	}
	return
}
