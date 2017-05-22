package core

import (
	"errors"
	"eye/integ"
	"fmt"
	"regexp"
	"gee/as"
)

type Eye struct {
	config       *Config
	accessFinder as.AccessFinder

	serviceFactory Factory
	checks         map[string]Check
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
}

func (o *Eye) reloadServiceFactory() {
	oldServiceFactory := o.serviceFactory
	o.serviceFactory = o.buildServiceFactory()
	if oldServiceFactory != nil {
		oldServiceFactory.Close()
	}

	o.checks = make(map[string]Check)

	//register checks
	o.registerMultiPing()
	o.registerValidateChecks()
	o.registerMultiValidates()
	o.registerCompares()
}

func (o *Eye) registerValidateChecks() {
	for _, item := range o.config.Validate {
		if len(item.Services) > 1 {
			for _, serviceName := range item.Services {
				o.registerValidateCheck(fmt.Sprintf("%v-%v", serviceName, item.Name),
					serviceName, item.Request)
			}
		} else if len(item.Services) > 0 {
			o.registerValidateCheck(item.Name,
				item.Services[0], item.Request)
		} else {
			Log.Info("No service defined for the check %v", item.Name)
		}
	}
}

func (o *Eye) registerValidateCheck(checkName string, serviceName string, request *QueryRequest) {
	check, err := o.buildCheck(serviceName, request)
	if err == nil {
		o.checks[checkName] = check
	} else {
		Log.Info("Can't build check '%v' because of '%v'", checkName, err)
	}
}

func (o *Eye) registerMultiPing() {
	for _, item := range o.config.PingAll {
		o.checks[item.Name] = &MultiPing{check: item, validator: o.PingAll}
	}

	for _, item := range o.config.PingAny {
		o.checks[item.Name] = &MultiPing{check: item, validator: o.PingAny}
	}
}

func (o *Eye) registerMultiValidates() {
	for _, item := range o.config.ValidateAll {
		o.checks[item.Name] = &MultiValidate{check: item, validator: o.ValidateAll}
	}

	for _, item := range o.config.ValidateAny {
		o.checks[item.Name] = &MultiValidate{check: item, validator: o.ValidateAny}
	}

	for _, item := range o.config.ValidateRunning {
		o.checks[item.Name] = &MultiValidate{check: item, validator: o.ValidateRunning}
	}
}

func (o *Eye) registerCompares() {
	for _, item := range o.config.CompareAll {
		if check, err := o.buildCompareCheck(item.Name, item.Services, false, item.Request,
			item.Request.compareValidator); err == nil {
			o.checks[item.Name] = check
		} else {
			item.logBuildCheckNotPossible(err)
		}
	}

	for _, item := range o.config.CompareAny {
		if check, err := o.buildCompareCheck(item.Name, item.Services, true, item.Request,
			item.Request.compareAnyValidator); err == nil {
			o.checks[item.Name] = check
		} else {
			item.logBuildCheckNotPossible(err)
		}
	}

	for _, item := range o.config.CompareRunning {
		if check, err := o.buildCompareCheck(item.Name, item.Services, true, item.Request,
			item.Request.compareValidator); err == nil {
			o.checks[item.Name] = check
		} else {
			item.logBuildCheckNotPossible(err)
		}
	}
}

func (o *Eye) buildServiceFactory() Factory {
	serviceFactory := NewFactory()
	for _, item := range o.config.MySql {
		serviceFactory.Add(&MySqlService{mysql: item, accessFinder: o.accessFinder})
	}

	for _, item := range o.config.Http {
		serviceFactory.Add(&HttpService{http: item, accessFinder: o.accessFinder})
	}

	for _, item := range o.config.Fs {
		serviceFactory.Add(&FsService{Fs: item})
	}

	for _, item := range o.config.Ps {
		serviceFactory.Add(&PsService{Ps: item})
	}
	return &serviceFactory
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

func (o *Eye) buildCheck(serviceName string, req *QueryRequest) (ret Check, err error) {
	var value interface{}
	value, err = o.liveChecks.GetOrBuild(req.CheckKey(serviceName), func() (interface{}, error) {
		service, err := o.serviceFactory.Find(serviceName)
		if err == nil {
			return service.NewСheck(req)
		} else {
			return nil, err
		}
	})
	if err == nil {
		ret = value.(Check)
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

func (o *Eye) getOrBuildCompareCheck(checkKey string, serviceNames []string, onlyRunning bool, req *CompareRequest,
	validator func(map[string]QueryResultInfo) error) (check Check, err error) {
	var value interface{}

	if value, err = o.liveChecks.GetOrBuild(checkKey, func() (interface{}, error) {
		return o.buildCompareCheck(checkKey, serviceNames, onlyRunning, req, validator)
	}); err == nil {
		check = value.(Check)
	}
	return
}

func (o *Eye) buildCompareCheck(checkKey string, serviceNames []string, onlyRunning bool, req *CompareRequest,
	validator func(map[string]QueryResultInfo) error) (check Check, err error) {

	var pattern *regexp.Regexp
	if len(req.QueryRequest.Expr) > 0 {
		pattern, err = regexp.Compile(req.QueryRequest.Expr)
		if err != nil {
			return
		}
	}

	checks := make([]Check, len(serviceNames))
	check = MultiCheck{info: checkKey, query: req.QueryRequest.Query, pattern: pattern, checks: checks,
		onlyRunning:     onlyRunning, validator: validator}

	var serviceCheck Check
	for i, serviceName := range serviceNames {
		serviceCheck, err = o.buildCheck(serviceName, req.QueryRequest)
		if err != nil {
			break
		}
		checks[i] = serviceCheck
	}
	return
}

func (o *Eye) query(serviceName string, req *QueryRequest) (data QueryResult, err error) {
	var buildCheck Check
	buildCheck, err = o.buildCheck(serviceName, req)
	if err == nil {
		data, err = buildCheck.Query()
	}
	return
}
