package core

import (
	"fmt"
)

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

func (o *Eye) registerValidateCheck(checkName string, serviceName string, request *ValidationRequest) {
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

func (o *Eye) getOrBuildCompareCheck(checkKey string, serviceNames []string, onlyRunning bool, req *ValidationRequest,
	validator func(map[string]QueryResultInfo) error) (check Check, err error) {
	var value interface{}

	if value, err = o.liveChecks.GetOrBuild(checkKey, func() (interface{}, error) {
		return o.buildCompareCheck(checkKey, serviceNames, onlyRunning, req, validator)
	}); err == nil {
		check = value.(Check)
	}
	return
}

func (o *Eye) buildCompareCheck(checkKey string, serviceNames []string, onlyRunning bool, req *ValidationRequest,
	validator func(map[string]QueryResultInfo) error) (check Check, err error) {

	checks := make([]Check, len(serviceNames))
	check = MultiCheck{info: checkKey, query: req.Query, checks: checks,
		onlyRunning:     onlyRunning, validator: validator}

	var serviceCheck Check
	for i, serviceName := range serviceNames {
		serviceCheck, err = o.buildCheck(serviceName, req)
		if err != nil {
			break
		}
		checks[i] = serviceCheck
	}
	return
}

func (o *Eye) query(serviceName string, req *ValidationRequest) (data QueryResults, err error) {
	var buildCheck Check
	buildCheck, err = o.buildCheck(serviceName, req)
	if err == nil {
		data, err = buildCheck.Query()
	}
	return
}

func (o *Eye) buildCheck(serviceName string, req *ValidationRequest) (ret Check, err error) {
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
