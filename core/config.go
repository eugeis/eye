package core

import (
	"encoding/json"
	"gee/cfg"
	"gee/lg"
	"path/filepath"
	"strings"
)

var Log = lg.NewLogger("EYE ")

type Config struct {
	Name         string `default:"Eye"`
	Port         int    `default:"3000"`
	Debug        bool   `default:true`
	ExportFolder string `default:"./export"`
	AppHome      string `default:"."`

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

	CompareRunning []*ValidateCheck
	CompareAll     []*ValidateCheck

	FieldsExporter []*FieldsExporter

	FileExporter []*FileExporter

	ConfigFiles    []string
	ConfigSuffixes []string
}

type PingCheck struct {
	Name     string
	Services []string
}

type ValidateCheck struct {
	Name     string
	Services []string
	Request  *ValidationRequest
}

type FieldsExporter struct {
	Name      string
	Query     string
	EvalExpr  string
	Fields    []string
	Separator string
	Services  []string
}

type FileExporter struct {
	Name           string
	Query          string
	EvalExpr       string
	Fields         []string
	SourceFileExpr string
	Services       []string
}

func LoadConfig(files []string, suffixes []string, appHome string) (ret *Config, err error) {
	ret = &Config{ConfigFiles: files, AppHome: appHome}
	err = cfg.Unmarshal(ret, files, suffixes)

	if err == nil {
		ret.calculateExportFolder()
		ret.Print()
	}
	return
}

func (o *Config) calculateExportFolder() (err error) {
	if len(o.ExportFolder) == 0 {
		o.ExportFolder = o.AppHome
	} else {
		if strings.HasPrefix(o.ExportFolder, "./") {
			o.ExportFolder = strings.Replace(o.ExportFolder, ".", o.AppHome+"/", 1)
		}
		o.ExportFolder, err = filepath.Abs(o.ExportFolder)
	}
	return
}

func (o *Config) Reload() (ret *Config, err error) {
	return LoadConfig(o.ConfigFiles, o.ConfigSuffixes, o.AppHome)
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
