package core

import (
	"errors"
	"regexp"
	"fmt"
	"bytes"
)

type QueryResultInfo struct {
	result QueryResult
	err    error
}

type MultiCheck struct {
	info        string
	query       string
	pattern     *regexp.Regexp
	checks      []Check
	onlyRunning bool
	validator   func(map[string]QueryResultInfo) error
}

func (o MultiCheck) Validate() error {
	return o.validator(o.checksData())
}

func (o MultiCheck) Query() (data QueryResult, err error) {
	return
}

func (o MultiCheck) Info() string {
	return o.info
}

func (o MultiCheck) checksData() (checksData map[string]QueryResultInfo) {
	checksData = make(map[string]QueryResultInfo)
	for _, check := range o.checks {
		data, err := check.Query()

		if err == nil {
			if o.pattern != nil {
				matches := o.pattern.FindSubmatch(data)
				if matches == nil {
					checksData[check.Info()] = QueryResultInfo{err: errors.New("No match")}
				} else if len(matches) > 0 {
					checksData[check.Info()] = QueryResultInfo{result: matches[1]}
					//Log.Debug("%v = %v", check.Info(), string(matches[1]))
				} else {
					checksData[check.Info()] = QueryResultInfo{result: matches[0]}
					//Log.Debug("%v = %v", check.Info(), string(matches[0]))
				}
			} else {
				checksData[check.Info()] = QueryResultInfo{result: data}
			}
		} else if !o.onlyRunning {
			checksData[check.Info()] = QueryResultInfo{result: data, err: err}
		}
	}
	return
}

type MultiPing struct {
	check     *PingCheck
	validator func([]string) error
}

func (o *MultiPing) Validate() (err error) {
	return o.validator(o.check.Services)
}

func (o MultiPing) Query() (data QueryResult, err error) {
	return
}

func (o MultiPing) Info() string {
	return o.check.Name
}

type MultiValidate struct {
	check     *ValidateCheck
	validator func([]string, *ValidationRequest) error
}

func (o *MultiValidate) Validate() (err error) {
	return o.validator(o.check.Services, o.check.Request)
}

func (o MultiValidate) Query() (data QueryResult, err error) {
	return
}

func (o MultiValidate) Info() string {
	return o.check.Name
}

func (o *CompareRequest) compareAnyValidator(checksData map[string]QueryResultInfo) (err error) {
	var firstName string
	var firstData QueryResult
	for serviceName, data := range checksData {
		if data.err != nil {
			if len(firstName) > 0 {
				firstData = data.result
				firstName = serviceName
			} else {
				matchError := match(firstData, data.result, o)
				if matchError != nil {
					err = errors.New(fmt.Sprintf("%v ~ %v => %v",
						firstName, serviceName, matchError))
				} else {
					err = nil
					break
				}
			}
		}
	}
	return
}

func (o *CompareRequest) compareValidator(checksData map[string]QueryResultInfo) error {
	return matchChecksData(checksData, o)
}

func (o *CompareCheck) logBuildCheckNotPossible(err error) {
	Log.Info("Can't build check '%v' because of '%v'", o.Name, err)
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
	if req.Tolerance > 0 {
		var x, y int
		var err error

		x, err = data1.Int()
		if err == nil {
			y, err = data2.Int()
		}
		if err == nil {
			diff := abs(x - y)
			if match := diff <= req.Tolerance; (!match && !req.ValidationRequest.Not) || (match && req.ValidationRequest.Not) {
				if req.ValidationRequest.Not {
					err = errors.New(fmt.Sprintf("abs(%v-%v)<=%v, less than tolerance %v", x, y, diff, req.Tolerance))
				} else {
					err = errors.New(fmt.Sprintf("abs(%v-%v)>%v, greater than tolerance %v", x, y, diff, req.Tolerance))
				}
			}
		} else {
			Log.Debug("Convertion to int not possible, use bytes comparison.")
			if match := bytes.Equal(data1, data2); (!match && !req.ValidationRequest.Not) || (match && req.ValidationRequest.Not) {
				if req.ValidationRequest.Not {
					err = errors.New("equal")
				} else {
					err = errors.New("not equal")
				}
			}
		}
	} else {
		if match := bytes.Equal(data1, data2); (!match && !req.ValidationRequest.Not) || (match && req.ValidationRequest.Not) {
			err = errors.New("not equal")
		}
	}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
