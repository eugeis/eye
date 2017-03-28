package core

import (
	"bytes"
	"errors"
	"eye/integ"
	"fmt"
	"regexp"
)

type Eye struct {
	config       *Config
	accessFinder AccessFinder

	serviceFactory Factory
	commandCache   integ.Cache
}

func NewEye(config *Config, accessFinder AccessFinder) (ret *Eye) {
	ret = &Eye{config: config, accessFinder: accessFinder, commandCache: integ.NewCache()}
	ret.reloadServiceFactory()
	return
}

func (o *Eye) Close() {
	if o.serviceFactory != nil {
		o.serviceFactory.Close()
		o.commandCache.Clear()
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
}

func (o *Eye) buildServiceFactory() Factory {
	serviceFactory := NewFactory()
	for _, item := range o.config.MySql {
		serviceFactory.Add(&MySqlService{mysql: item, accessFinder: o.accessFinder})
	}

	for _, item := range o.config.Http {
		serviceFactory.Add(&HttpService{http: item, accessFinder: o.accessFinder})
	}
	return &serviceFactory
}

func (o *Eye) Ping(serviceName string) (err error) {
	service, err := o.serviceFactory.Find(serviceName)
	if err == nil {
		err = service.Ping()
	}
	return
}

func (o *Eye) Check(checkName string) (err error) {
	return
}

func (o *Eye) Validate(serviceName string, req *QueryRequest) (err error) {
	if req.Query == "" {
		log.Debug(fmt.Sprintf("ping instead of validator, because no query defined for %v", serviceName))
		return o.Ping(serviceName)
	}

	var check Check
	check, err = o.buildCheck(serviceName, req)
	if err == nil {
		err = check.Validate()
	}
	return
}

func (o *Eye) buildCheck(serviceName string, req *QueryRequest) (ret Check, err error) {
	var value interface{}
	value, err = o.commandCache.GetOrBuild(req.CheckKey(serviceName), func() (interface{}, error) {
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
	checkKey := req.ChecksKey("any", serviceNames)
	check, err = o.buildCompareCheck(checkKey, serviceNames, true, req, func(checksData map[string]QueryResultInfo) (checkError error) {
		var firstName string
		var firstData QueryResult
		for serviceName, data := range checksData {
			if data.err != nil {
				if len(firstName) > 0 {
					firstData = data.result
					firstName = serviceName
				} else {
					matchError := match(firstData, data.result, req)
					if matchError != nil {
						checkError = errors.New(fmt.Sprintf("%v ~ %v => %v",
							firstName, serviceName, matchError))
					} else {
						checkError = nil
						break
					}
				}
			}
		}
		return
	})

	if err == nil {
		err = check.Validate()
	}
	return
}

func (o *Eye) CompareRunning(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	checkKey := req.ChecksKey("running", serviceNames)
	check, err = o.buildCompareCheck(checkKey, serviceNames, true, req, func(checksData map[string]QueryResultInfo) error {
		return matchChecksData(checksData, req)
	})

	if err == nil {
		err = check.Validate()
	}
	return
}

func (o *Eye) CompareAll(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	checkKey := req.ChecksKey("all", serviceNames)
	check, err = o.buildCompareCheck(checkKey, serviceNames, false, req, func(checksData map[string]QueryResultInfo) error {
		return matchChecksData(checksData, req)
	})

	if err == nil {
		err = check.Validate()
	}
	return
}

func matchChecksData(checksData map[string]QueryResultInfo, req *CompareRequest) (err error) {
	var firstName string
	var firstData QueryResult
	for serviceName, queryResult := range checksData {
		if queryResult.err == nil {
			if len(firstName) == 0 {
				firstData = queryResult.result
				firstName = serviceName
			} else {
				matchError := match(firstData, queryResult.result, req)
				if matchError != nil {
					err = errors.New(fmt.Sprintf("%v ~ %v failed: %v",
						firstName, serviceName, matchError))
					break
				}
			}
		} else {
			err = errors.New(fmt.Sprintf("%v failed: %v", serviceName, queryResult.err.Error()))
		}
	}
	return
}

func match(data1 QueryResult, data2 QueryResult, req *CompareRequest) (err error) {
	var matchErr error
	if req.Tolerance > 0 {
		var x, y int
		var err error

		x, err = data1.Int()
		if err == nil {
			y, err = data2.Int()
		}
		if err == nil {
			diff := abs(x - y)
			if diff > req.Tolerance {
				matchErr = errors.New(fmt.Sprintf("abs(%v-%v)=%v, greater than tolerance %v", x, y, diff, req.Tolerance))
			}
		} else {
			log.Debug("Convertion to int not possible, use bytes comparison.")
			if !bytes.Equal(data1, data2) {
				matchErr = errors.New("not equal")
			}
		}
	} else {
		if !bytes.Equal(data1, data2) {
			matchErr = errors.New("not equal")
		}
	}
	if !req.Not {
		err = matchErr
	} else if matchErr == nil {
		err = errors.New("compare request matches but must not")
	}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (o *Eye) buildCompareCheck(checkKey string, serviceNames []string, onlyRunning bool, req *CompareRequest,
	validator func(map[string]QueryResultInfo) error) (check Check, err error) {
	var value interface{}

	value, err = o.commandCache.GetOrBuild(checkKey, func() (check interface{}, err error) {
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
	})

	if err == nil {
		check = value.(Check)
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
