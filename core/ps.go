package core

import (
	"github.com/shirou/gopsutil/process"
	"eye/integ"
	"runtime"
	"github.com/StackExchange/wmi"
	"gopkg.in/Knetic/govaluate.v2"
)

type Ps struct {
	Name        string
	PingRequest *ValidationRequest
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
			o.pingCheck, err = o.newСheck(&ValidationRequest{})
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

func (o *PsService) NewСheck(req *ValidationRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *PsService) newСheck(req *ValidationRequest) (ret *PsCheck, err error) {
	var eval *govaluate.EvaluableExpression
	if eval, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	ret = &PsCheck{info: req.CheckKey("Ps"), service: o, eval: eval, all: req.All}
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
	all     bool
	eval    *govaluate.EvaluableExpression
}

func (o PsCheck) Info() string {
	return o.info
}

func (o PsCheck) Validate() (err error) {
	return validate(o, o.eval, o.all)
}

func (o PsCheck) Query() (ret QueryResults, err error) {
	var items []*Proc
	if items, err = o.service.Processes(); err == nil {
		ret = make([]QueryResult, len(items))
		for i, item := range items {
			ret[i] = &MapQueryResult{item.ToMap()}
		}
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

func (o *Proc) ToMap() (ret map[string]interface{}) {
	return map[string]interface{}{
		"Id": o.Id, "Name": o.Name, "Status": o.Status, "Cmdline": o.Cmdline, "Path": o.Path}
}
