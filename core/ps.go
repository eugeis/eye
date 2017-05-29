package core

import (
	"regexp"
	"github.com/shirou/gopsutil/process"
	"encoding/json"
	"eye/integ"
	"runtime"
	"github.com/StackExchange/wmi"
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
	if err = o.Init(); err == nil {
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
			query:       pattern[0], pattern: pattern[1], not: req.Not}
	}
	return
}

func (o *PsService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
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
	if runtime.GOOS == "windows" {
		var dst []process.Win32_Process
		q := wmi.CreateQuery(&dst, "")
		if err = wmi.Query(q, &dst); err == nil {
			ret = make([]*Proc, len(dst))
			for i, wp := range dst {
				proc := &Proc{Id: int(wp.ProcessID)}
				ret[i] = proc

				proc.Name = wp.Name
				proc.Status = *wp.Status
				proc.Cmdline = *wp.CommandLine
				proc.Path = *wp.ExecutablePath
			}
		}
	} else {
		var pids []int32
		if pids, err = process.Pids(); err == nil {
			ret = make([]*Proc, len(pids))

			for i, pid := range pids {
				proc := &Proc{Id: int(pid)}
				ret[i] = proc

				p, _ := process.NewProcess(int32(pid))

				proc.Name, _ = p.Name()
				proc.Status, _ = p.Status()
				proc.Cmdline, _ = p.Cmdline()
			}

		}
	}
	return
}

//buildCheck
type PsCheck struct {
	info    string
	service *PsService
	not     bool
	query   *regexp.Regexp
	pattern *regexp.Regexp
}

func (o PsCheck) Info() string {
	return o.info
}

func (o PsCheck) Validate() (err error) {
	return validate(o, o.pattern, o.not)
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

type Proc struct {
	Id      int
	Name    string
	Status  string
	Cmdline string
	Path    string
}
