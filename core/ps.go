package core

import (
	"errors"
	"fmt"
	"regexp"
	"github.com/shirou/gopsutil/process"
)

type Ps struct {
	Name        string
	PingRequest *QueryRequest
}

type PsService struct {
	Ps        *Ps
	pingCheck *PsCheck
}

func (o *PsService) Name() string {
	return o.Ps.Name
}

func (o *PsService) Init() (err error) {
	if o.pingCheck == nil {
		o.pingCheck, err = o.newСheck(o.Ps.PingRequest)
		if err != nil {
			o.Close()
		}
	}
	return
}

func (o *PsService) Close() {
	o.pingCheck = nil
}

func (o *PsService) Ping() (err error) {
	if err = o.Init(); err != nil {
		err = o.pingCheck.Validate()
	}
	return
}

func (o *PsService) NewСheck(req *QueryRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *PsService) newСheck(req *QueryRequest) (ret *PsCheck, err error) {
	return
}

//buildCheck
type PsCheck struct {
	info    string

	process *process.Process

	pattern *regexp.Regexp
	service *PsService
}

func (o PsCheck) Info() string {
	return o.info
}

func (o PsCheck) Validate() (err error) {
	data, err := o.Query()

	if err == nil && o.pattern != nil {
		if !o.pattern.Match(data) {
			err = errors.New(fmt.Sprintf("No match for %v", o.info))
		}
	}
	return
}

func (o PsCheck) Query() (data QueryResult, err error) {
	//p,_ := process.NewProcess(1)
	return
}
