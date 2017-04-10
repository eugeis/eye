package core

import (
	"errors"
	"fmt"
	"regexp"
	"os"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
)

type Fs struct {
	Name        string
	File        string
	PingRequest *QueryRequest
}

type FsService struct {
	Fs        *Fs
	pingCheck *FsCheck
}

func (o *FsService) Name() string {
	return o.Fs.Name
}

func (o *FsService) Init() (err error) {
	if o.pingCheck == nil {
		if o.Fs.PingRequest != nil {
			o.pingCheck, err = o.newСheck(o.Fs.PingRequest)
		} else {
			o.pingCheck, err = o.newСheck(nil)
		}
		if err != nil {
			o.Close()
		}
	}
	return
}

func (o *FsService) Close() {
	o.pingCheck = nil
}

func (o *FsService) Ping() (err error) {
	if err = o.Init(); err != nil {
		err = o.pingCheck.Validate()
	}
	return
}

func (o *FsService) NewСheck(req *QueryRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *FsService) newСheck(req *QueryRequest) (ret *FsCheck, err error) {
	var pattern *regexp.Regexp
	if req != nil {
		if pattern, err = compilePattern(req.Expr); err == nil {
			if req.Query != "" {
				ret = &FsCheck{info: req.CheckKey("fs"), file: filepath.Join(o.Fs.File, req.Query), pattern: pattern}
			} else {
				ret = &FsCheck{info: req.CheckKey("fs"), file: o.Fs.File, pattern: pattern}
			}
		}
	} else {
		ret = &FsCheck{info: req.CheckKey("fs"), file: o.Fs.File}
	}
	return
}

//buildCheck
type FsCheck struct {
	info    string
	file    string
	pattern *regexp.Regexp
}

func (o FsCheck) Info() string {
	return o.info
}

func (o FsCheck) Validate() (err error) {
	data, err := o.Query()

	if err == nil && o.pattern != nil {
		if !o.pattern.Match(data) {
			err = errors.New(fmt.Sprintf("No match for %v", o.info))
		}
	}
	return
}

func (o FsCheck) Query() (data QueryResult, err error) {
	var file os.FileInfo

	file, err = os.Stat(o.file)

	if file.IsDir() {
		files, _ := ioutil.ReadDir(o.file)
		data, err = json.Marshal(files)
	} else {
		data, err = json.Marshal(file)
	}
	return
}
