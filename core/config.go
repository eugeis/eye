package core

import (
	"eye/integ"
	"fmt"
	"github.com/jinzhu/configor"
	"strings"
)

var log = integ.Log

type Config struct {
	Name  string `default:"Eye"`
	Port  int    `default:"3000"`
	Debug bool   `default:true`
	Token string

	MySql []*MySql
	Http  []*Http

	configFiles []string
}

func Load(fileNames []string) (ret *Config, err error) {
	ret = &Config{}
	err = configor.Load(ret, fileNames...)

	ret.Print()

	//ignore, https://github.com/jinzhu/configor/issues/6
	if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
		err = nil
	}
	return
}

func (o *Config) Reload() (ret *Config, err error) {
	return Load(o.configFiles)
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
