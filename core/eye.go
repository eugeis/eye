package core

import (
	"errors"
	"eye/integ"
	"fmt"
	"gee/as"
	"os"
)

type Eye struct {
	config       *Config
	accessFinder as.AccessFinder

	serviceFactory Factory
	checks         map[string]Check
	exporters      map[string]Exporter
	liveChecks     integ.Cache
}

func NewEye(config *Config, accessFinder as.AccessFinder) (ret *Eye) {
	ret = &Eye{config: config, accessFinder: accessFinder, liveChecks: integ.NewCache()}
	ret.reloadServiceFactory()
	return
}

func (o *Eye) Close() {
	if o.serviceFactory != nil {
		o.serviceFactory.Close()
		o.liveChecks.Clear()
	}
}

func (o *Eye) UpdateConfig(config *Config) {
	o.Close()
	o.config = config
	o.reloadServiceFactory()
	if len(o.config.ExportFolder) > 0 {
		os.MkdirAll(o.config.ExportFolder, 0777)
	}
}

func (o *Eye) Ping(serviceName string) (err error) {
	var service Service
	if service, err = o.serviceFactory.Find(serviceName); err == nil {
		err = service.Ping()
	}
	return
}

func (o *Eye) Check(checkName string) (err error) {
	if check, ok := o.checks[checkName]; ok {
		err = check.Validate()
	} else {
		err = errors.New(fmt.Sprintf("There is no check '%v' available", checkName))
	}
	return
}

func (o *Eye) Export(exportName string, params map[string]string) (err error) {
	if exporter, ok := o.exporters[exportName]; ok {
		err = exporter.Export(params)
	} else {
		err = errors.New(fmt.Sprintf("There is no exporter '%v' available", exportName))
	}
	return
}

func (o *Eye) Validate(serviceName string, req *QueryRequest) (err error) {
	if req.Query == "" {
		Log.Debug(fmt.Sprintf("ping instead of validator, because no query defined for %v", serviceName))
		return o.Ping(serviceName)
	}

	var check Check
	if check, err = o.buildCheck(serviceName, req); err == nil {
		err = check.Validate()
	}
	return
}

func (o *Eye) PingAny(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err == nil {
			break
		}
	}
	return
}

func (o *Eye) PingAll(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err != nil {
			break
		}
	}
	return
}

func (o *Eye) ValidateAny(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err == nil {
			break
		}
	}
	return
}

func (o *Eye) ValidateRunning(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		errTemp := o.Ping(serviceName)
		if errTemp == nil {
			err = o.Validate(serviceName, req)
			if err != nil {
				break
			}
		}

	}
	return
}

func (o *Eye) ValidateAll(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err != nil {
			break
		}
	}
	return
}

func (o *Eye) CompareAny(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	if check, err = o.getOrBuildCompareCheck(req.ChecksKey("any", serviceNames), serviceNames,
		true, req, req.compareAnyValidator); err == nil {
		err = check.Validate()
	}
	return
}

func (o *Eye) CompareRunning(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	check, err = o.getOrBuildCompareCheck(req.ChecksKey("running", serviceNames), serviceNames,
		true, req, req.compareValidator)

	if err == nil {
		err = check.Validate()
	}
	return
}

func (o *Eye) CompareAll(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	if check, err = o.getOrBuildCompareCheck(req.ChecksKey("all", serviceNames), serviceNames,
		false, req, req.compareValidator); err == nil {
		err = check.Validate()
	}
	return
}
