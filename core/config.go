package core

import (
	"eye/integ"
	"github.com/jinzhu/configor"
	"strings"
	"encoding/json"
	"os"
)

var l = integ.Log

type Config struct {
	Name  string `default:"Eye"`
	Port  int    `default:"3000"`
	Debug bool   `default:true`

	MySql []*MySql
	Http  []*Http

	ConfigFile string
}

func LoadConfig(file string) (ret *Config, err error) {
	if _, err = os.Stat(file); err == nil {
		ret = &Config{ConfigFile: file}
		err = configor.Load(ret, file)

		ret.Print()

		//ignore, https://github.com/jinzhu/configor/issues/6
		if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
			err = nil
		}
	}
	return
}

func (o *Config) Reload() (ret *Config, err error) {
	return LoadConfig(o.ConfigFile)
}

func (o *Config) Print() {
	json, err := json.MarshalIndent(o, "", "\t")
	if err == nil {
		l.Info(string(json))
	}
}

func fillAccessData(security *Security, file string) (err error) {
	err = configor.Load(security, file)
	//ignore, https://github.com/jinzhu/configor/issues/6
	if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
		err = nil
	}
	return
}
