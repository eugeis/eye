package core

import (
	"rest/integ"
	"fmt"
	"bytes"
	"regexp"
	"errors"
)

type Operator interface {
}

type Controller struct {
	config         *Config
	serviceFactory Factory
	commandCache   integ.Cache
}

func NewController(config *Config) Controller {
	return Controller{config: config, serviceFactory: config.ServiceFactory(), commandCache: integ.NewCache()}
}

func (o Controller) Close() {
	if o.serviceFactory != nil {
		o.serviceFactory.Close()
	}
	o.commandCache = nil
	o.config = nil
	o.serviceFactory = nil
}

func (o Controller) Ping(serviceName string) (err error) {
	service, err := o.serviceFactory.Find(serviceName)
	if err == nil {
		err = service.Ping()
	}
	return
}

func (o Controller) Check(checkName string) (err error) {
	return
}

func (o Controller) Validate(serviceName string, req *QueryRequest) (err error) {
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

func (o Controller) buildCheck(serviceName string, req *QueryRequest) (ret Check, err error) {
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

func (o Controller) PingAny(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err == nil {
			break
		}
	}
	return
}

func (o Controller) PingAll(serviceNames []string) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Ping(serviceName)
		if err != nil {
			break
		}
	}
	return
}

func (o Controller) ValidateAny(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err == nil {
			break
		}
	}
	return
}

func (o Controller) ValidateRunning(serviceNames []string, req *QueryRequest) (err error) {
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

func (o Controller) ValidateAll(serviceNames []string, req *QueryRequest) (err error) {
	for _, serviceName := range serviceNames {
		err = o.Validate(serviceName, req)
		if err != nil {
			break
		}
	}
	return
}

func (o Controller) CompareAny(serviceNames []string, req *CompareRequest) (err error) {
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

func (o Controller) CompareRunning(serviceNames []string, req *CompareRequest) (err error) {
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

func (o Controller) CompareAll(serviceNames []string, req *CompareRequest) (err error) {
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
						break
					}
				}
			} else {
				checkError = errors.New(fmt.Sprintf("%s != $s", firstName, serviceName))
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

func (o Controller) compareCheck(serviceNames []string, req *CompareRequest, validator func(map[string][]byte) error) (check Check, err error) {
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
		check = multiChecks{info: checkKey, query: req.QueryRequest.Query, pattern: pattern, checks: checks, validator: validator}

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

func (o Controller) query(serviceName string, req *QueryRequest) (data []byte, err error) {
	//var buildCheck Check
	//buildCheck, err = o.buildCheck(serviceName, req)
	if err == nil {
		//data, err = buildCheck.Query()
	}
	return
}

//buildCheck
type multiChecks struct {
	info      string
	query     string
	pattern   *regexp.Regexp
	checks    []Check
	validator func(map[string][]byte) error
}

func (o multiChecks) Validate() error {
	return o.validator(o.checksData())
}

func (o multiChecks) Query() (data []byte, err error) {
	return
}

func (o multiChecks) Info() string {
	return o.info
}

func (o multiChecks) checksData() (checksData map[string][]byte) {
	checksData = make(map[string][]byte)
	for _, check := range o.checks {
		data, err := check.Query()
		if err == nil {
			if o.pattern != nil {
				matches := o.pattern.FindSubmatch(data)
				if matches == nil {
					checksData[check.Info()] = nil
				} else if len(matches) > 0 {
					checksData[check.Info()] = matches[1]
				} else {
					checksData[check.Info()] = matches[0]
				}
			} else {
				checksData[check.Info()] = data
			}
		} else {
			checksData[check.Info()] = nil
		}
	}
	return
}
