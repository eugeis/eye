package core

import (
	"encoding/json"
	"gee/cfg"
	"gee/lg"
)

var Log = lg.NewLogger("EYE ")

type Config struct {
	Name         string `default:"Eye"`
	Port         int    `default:"3000"`
	Debug        bool   `default:true`
	ExportFolder string `default:"./export"`

	MySql   []*MySql
	Http    []*Http
	Fs      []*Fs
	Ps      []*Ps
	Elastic []*Elastic

	PingAny []*PingCheck
	PingAll []*PingCheck

	Validate []*ValidateCheck

	ValidateAny     []*ValidateCheck
	ValidateRunning []*ValidateCheck
	ValidateAll     []*ValidateCheck

	CompareAny     []*CompareCheck
	CompareRunning []*CompareCheck
	CompareAll     []*CompareCheck

	FieldsExporter []*FieldsExporter

	ConfigFiles    []string
	ConfigSuffixes []string
}

type PingCheck struct {
	Name     string
	Services []string
}

type ValidateCheck struct {
	Name     string
	Request  *QueryRequest
	Services []string
}

type CompareCheck struct {
	Name     string
	Request  *CompareRequest
	Services []string
}

type FieldsExporter struct {
	Name      string
	Query     string
	Fields    []string
	Separator string
	Services  []string
}

func LoadConfig(files []string, suffixes []string) (ret *Config, err error) {
	ret = &Config{ConfigFiles: files}
	err = cfg.Unmarshal(ret, files, suffixes)

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
		Log.Info(string(json))
	}
}

func (o *Config) ExtractAccessKeys() (ret []string) {
	ret = make([]string, len(o.MySql)+len(o.Http))
	for i, item := range o.MySql {
		ret[i] = item.AccessKey
	}
	var pre = len(o.MySql)
	for i, item := range o.Http {
		ret[pre+i] = item.AccessKey
	}
	return
}
