package core

import (
	"errors"
	"fmt"
	"regexp"
	"github.com/shirou/gopsutil/process"
	"encoding/json"
	"github.com/shirou/gopsutil/net"
	"eye/integ"
)

type Ps struct {
	Name        string
	PingRequest *QueryRequest
}

type PsService struct {
	Ps        *Ps
	pingCheck *PsCheck

	procs *integ.TimeoutObjectCache
}

func (o *PsService) Name() string {
	return o.Ps.Name
}

func (o *PsService) Init() (err error) {
	if o.pingCheck == nil {
		if o.Ps.PingRequest != nil {
			o.pingCheck, err = o.newСheck(o.Ps.PingRequest)
		} else {
			o.pingCheck, err = o.newСheck(&QueryRequest{})
		}
		if err != nil {
			o.Close()
		} else {
			o.procs = integ.NewObjectCache(func() (interface{}, error) { return o.buildProcesses() })
		}

	}
	return
}

func (o *PsService) Close() {
	o.pingCheck = nil
	o.procs = nil
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
	var pattern []*regexp.Regexp

	pattern, err = compilePatterns(req.Query, req.Expr)
	if err == nil {
		ret = &PsCheck{info: req.CheckKey("Ps"), service: o,
			query:       pattern[0], pattern: pattern[1]}
	}
	return
}

func (o PsService) Processes() (ret []*Proc, err error) {
	value, err := o.procs.Get()
	if err == nil {
		ret = value.([]*Proc)
	}
	return
}

func (o PsService) buildProcesses() (ret []*Proc, err error) {
	var pids []int32
	if pids, err = process.Pids(); err == nil {
		ret = make([]*Proc, len(pids))

		for i, pid := range pids {
			proc := &Proc{Id: pid}
			ret[i] = proc

			p, _ := process.NewProcess(int32(pid))

			proc.Name, _ = p.Name()
			proc.Status, _ = p.Status()
			proc.Cmdline, _ = p.Cmdline()
			//proc.Connections, _ = p.Connections()
		}
	}
	return
}

//buildCheck
type PsCheck struct {
	info    string
	service *PsService
	query   *regexp.Regexp
	pattern *regexp.Regexp
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

func (o PsCheck) Query() (ret QueryResult, err error) {
	var procs []*Proc
	if o.query != nil {
		procs, err = o.service.Processes()
	} else {
		procs, err = o.service.Processes()
	}
	if err == nil {
		ret, err = json.Marshal(procs)
	}
	return
}

func (o PsCheck) queryToMap() (ret []*Proc, err error) {
	var pids []int32
	if pids, err = process.Pids(); err == nil {
		ret = make([]*Proc, len(pids))

		for i, pid := range pids {
			proc := &Proc{Id: pid}
			ret[i] = proc

			p, _ := process.NewProcess(int32(pid))

			proc.Name, _ = p.Name()
			//proc.Status, _ = p.Status()
			//proc.Cmdline, _ = p.Cmdline()
			//proc.Connections, _ = p.Connections()
		}
	}
	return
}

type Proc struct {
	Id          int32
	Name        string
	Status      string
	Cmdline     string
	Connections []net.ConnectionStat
}
