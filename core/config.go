package core

import (
	"eye/integ"
	"github.com/jinzhu/configor"
	"strings"
	"encoding/json"
	"os"
	"regexp"
	"bytes"
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

	ConfigFile string
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

func LoadConfig(file string) (ret *Config, err error) {
	if _, err = os.Stat(file); err == nil {
		ret = &Config{ConfigFile: file}
		err = configor.Load(ret, file)

		//ignore, https://github.com/jinzhu/configor/issues/6
		if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
			err = nil
		}

		if err == nil {
			//ret.replaceEnvironmentVariables()
			ret.Print()
		}
	}
	return
}

func (o *Config) replaceEnvironmentVariables() (err error) {
	pattern := "(.+?)\\$\\{(.+?)\\}(.+?)"
	regexp := regexp.MustCompile(pattern)
	o.Name = replaceEnvVariables(o.Name, regexp)
	return
}

func replaceEnvVariables(string string, regexp *regexp.Regexp) (ret string) {
	if matches := regexp.FindAllStringSubmatch(string, -1); matches != nil {
		var str bytes.Buffer
		for _, group := range matches {
			str.WriteString(group[1])
			envKey := group[2]
			if envValue := os.Getenv(envKey); len(envValue) > 0 {
				str.WriteString(envValue)
			} else {
				l.Info("No environment value found for %v", envKey)
				str.WriteString(envKey)
			}
			str.WriteString(group[3])
			l.Info("%v", group)
		}
		ret = str.String()
	} else {
		ret = string
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
