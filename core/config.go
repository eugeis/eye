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

	SecurityType SecurityType

	MySql []*MySql
	Http  []*Http

	ConfigFile string
}

type SecurityType struct {
	Type string `default:"file"` //file, console, vault

	//file
	File string `default:"security.yml"`

	//vault
	VaultAddress string
	Token        string
}

func Load(configFile string) (ret *Config, err error) {
	ret = &Config{ConfigFile: configFile}
	err = configor.Load(ret, configFile)

	ret.Print()

	//ignore, https://github.com/jinzhu/configor/issues/6
	if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
		err = nil
	}
	return
}

func (o *Config) Reload() (ret *Config, err error) {
	return Load(o.ConfigFile)
}

func (o *Config) Print() {
	json, err := json.MarshalIndent(o, "", "\t")
	if err == nil {
		log.Info(string(json))
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
