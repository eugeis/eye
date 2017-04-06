package core

import (
	"eye/integ"
	"github.com/jinzhu/configor"
	"strings"
	"encoding/json"
	"eye/conf"
)

var l = integ.Log

type Config struct {
	Name  string `default:"Eye"`
	Port  int    `default:"3000"`
	Debug bool   `default:true`

	MySql []*MySql
	Http  []*Http

	Validate []*ValidateCheck

	ValidateAny     []*ValidateCheck
	ValidateRunning []*ValidateCheck
	ValidateAll     []*ValidateCheck

	CompareAny     []*CompareCheck
	CompareRunning []*CompareCheck
	CompareAll     []*CompareCheck

	ConfigFiles    []string
	ConfigSuffixes []string
}

type ValidateCheck struct {
	Name     string
	Services []string
	Request  *QueryRequest
}

type CompareCheck struct {
	Name     string
	Services []string
	Request  *CompareRequest
}

func LoadConfig(files []string, suffixes []string) (ret *Config, err error) {
	ret = &Config{ConfigFiles: files}
	err = conf.Unmarshal(ret, files, suffixes)

	if err == nil {
		ret.Print()
	}

	return
}

func (o *Config) Reload() (ret *Config, err error) {
	return LoadConfig(o.ConfigFiles, o.ConfigSuffixes)
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
