package core

import (
	"rest/integ"
	"fmt"
	"bytes"
	"regexp"
	"errors"
	"strconv"
)

type Operator interface {
}

type Controller struct {
	Config         *Config
	serviceFactory Factory
	commandCache   integ.Cache
	configFiles    string
}

func NewController(configFiles string) (ret *Controller, err error) {
	ret = &Controller{configFiles: configFiles, commandCache: integ.NewCache()}
	err = ret.ReloadConfig()
	return
}

func (o *Controller) Close() {
	if o.serviceFactory != nil {
		o.serviceFactory.Close()
		o.commandCache.Clear()
	}
}

func (o *Controller) ReloadConfig() (err error) {
	o.Close()
	o.Config, err = Load(o.configFiles)
	if err == nil {
		o.Config.Print()
		o.reloadServiceFactory()
	}
	return
}

func (o *Controller) reloadServiceFactory() {
	oldServiceFactory := o.serviceFactory
	o.serviceFactory = o.buildServiceFactory()
	if oldServiceFactory != nil {
		oldServiceFactory.Close()
	}
}

func (o *Controller) buildServiceFactory() Factory {
	serviceFactory := NewFactory()
	for _, service := range o.Config.MySql {
		serviceFactory.Add(service)
	}

	for _, service := range o.Config.Http {
		serviceFactory.Add(service)
	}
	return &serviceFactory
}

func (o *Controller) Ping(serviceName string) (err error) {
	service, err := o.serviceFactory.Find(serviceName)
	if err == nil {
		err = service.Ping()
	}
	return
}

func (o *Controller) Check(checkName string) (err error) {
	return
}

func (o *Controller) Validate(serviceName string, req *QueryRequest) (err error) {
	if req.Query == "" {
		log.Debug(fmt.Sprintf("ping instead of validator, because no query defined for %s", serviceName))
		return o.Ping(serviceName)
	}

	var check Check
	check, err = o.buildCheck(serviceName, req)
	if err == nil {
		err = check.Validate()
	}
	return
}

func (o *Controller) buildCheck(serviceName string, req *QueryRequest) (ret Check, err error) {
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

func (o *Controller) PingAny(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err == nil {
			break
		}
	}
	return
}

func (o *Controller) PingAll(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err != nil {
			break
		}
	}
	return
}

func (o *Controller) ValidateAny(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err == nil {
			break
		}
	}
	return
}

func (o *Controller) ValidateRunning(serviceNames []string, req *QueryRequest) (err error) {
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

func (o *Controller) ValidateAll(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err != nil {
			break
		}
	}
	return
}

func (o *Controller) CompareAny(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	check, err = o.compareCheck(serviceNames, req, func(checksData map[string][]byte) (checkError error) {
		var firstName string
		var firstData []byte
		for serviceName, data := range checksData {
			if data != nil {
				if firstData == nil {
					firstData = data
					firstName = serviceName
				} else {
					if !bytes.Equal(firstData, data) {
						checkError = errors.New(fmt.Sprintf("%s != $s", firstName, serviceName))
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

func (o *Controller) CompareRunning(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	check, err = o.compareCheck(serviceNames, req, func(checksData map[string][]byte) (checkError error) {
		var firstName string
		var firstData []byte
		for serviceName, data := range checksData {
			//check only running
			if data != nil {
				if firstData == nil {
					firstData = data
					firstName = serviceName
				} else {
					if !bytes.Equal(firstData, data) {
						checkError = errors.New(fmt.Sprintf("%s != $s", firstName, serviceName))
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

func (o *Controller) CompareAll(serviceNames []string, req *CompareRequest) (err error) {
	var check Check
	check, err = o.compareCheck(serviceNames, req, func(checksData map[string][]byte) (checkError error) {
		var firstName string
		var firstData []byte

		for serviceName, data := range checksData {
			if data != nil {
				if firstData == nil {
					firstData = data
					firstName = serviceName
				} else {
					if !match(firstData, data, req) {
						checkError = errors.New(fmt.Sprintf("%s ~ $s", firstName, serviceName))
						break
					}
				}
			} else {
				checkError = errors.New(fmt.Sprintf("%s ~ %s", firstName, serviceName))
				break
			}
		}
		return
	})

	if err == nil {
		err = check.Validate()
	}
	return
}

func match(data1 []byte, data2 []byte, req *CompareRequest) (ret bool) {
	match := false
	if req.Tolerance > 0 {
		var x, y int
		var err error
		s1 := string(data1)
		s2 := string(data2)

		x, err = strconv.Atoi(s1)
		if err == nil {
			y, err = strconv.Atoi(s2)
		}
		if err == nil {
			match = abs(x-y) <= req.Tolerance
		} else {
			log.Debug("Covert of %s and %s to int not possible, use bytes comparison.")
			match = bytes.Equal(data1, data2)
		}
	} else {
		match = bytes.Equal(data1, data2)
	}
	return (!req.Not && match) || (req.Not && !match)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (o *Controller) compareCheck(serviceNames []string, req *CompareRequest, validator func(map[string][]byte) error) (check Check, err error) {
	var value interface{}

	checkKey := req.QueryRequest.ChecksKey(serviceNames)
	value, err = o.commandCache.GetOrBuild(checkKey, func() (check interface{}, err error) {
		var pattern *regexp.Regexp
		if len(req.QueryRequest.Expr) > 0 {
			pattern, err = regexp.Compile(req.QueryRequest.Expr)
			if err != nil {
				return
			}
		}

		checks := make([]Check, len(serviceNames))
		check = MultiCheck{info: checkKey, query: req.QueryRequest.Query, pattern: pattern, checks: checks, validator: validator}

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
		err = check.Validate()
	}
	return
}

func (o *Controller) query(serviceName string, req *QueryRequest) (data []byte, err error) {
	//var buildCheck Check
	//buildCheck, err = o.buildCheck(serviceName, req)
	if err == nil {
		//data, err = buildCheck.Query()
	}
	return
}
