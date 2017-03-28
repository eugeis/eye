package core

import (
	"eye/integ"
	"github.com/jinzhu/configor"
	"strings"
	"encoding/json"
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

func Load(configFiles []string) (ret *Config, err error) {
	ret = &Config{configFiles: configFiles}
	err = configor.Load(ret, configFiles...)

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

func (o *Config) Print() {
	json, err := json.MarshalIndent(o, "", "\t")
	if err == nil {
		log.Info(string(json))
	}
}
