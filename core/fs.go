package core

import (
	"errors"
	"fmt"
	"regexp"
	"os"
)

type Fs struct {
	Name        string
	AccessKey   string
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
	if o.pingCheck == nil && o.Fs.PingRequest != nil {
		o.pingCheck, err = o.newСheck(o.Fs.PingRequest)
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
		if o.pingCheck != nil {
			err = o.pingCheck.Validate()
		} else {
			_, err = os.Stat(o.Fs.File)
		}
	}
	return
}

func (o *FsService) NewСheck(req *QueryRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *FsService) newСheck(req *QueryRequest) (ret *FsCheck, err error) {
	return
}

//buildCheck
type FsCheck struct {
	info    string
	pattern *regexp.Regexp
	service *FsService
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
	return
}
