package core

import (
	"errors"
	"fmt"
	"gopkg.in/Knetic/govaluate.v2"
)

type QueryResultInfo struct {
	result QueryResults
	err    error
}

type MultiCheck struct {
	info        string
	query       string
	checks      []Check
	onlyRunning bool
	validator   func(map[string]QueryResultInfo) error
}

func (o MultiCheck) Validate() error {
	return o.validator(o.checksData())
}

func (o MultiCheck) Query() (data QueryResults, err error) {
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
			checksData[check.Info()] = QueryResultInfo{result: data}
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

func (o MultiPing) Query() (data QueryResults, err error) {
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

func (o MultiValidate) Query() (data QueryResults, err error) {
	return
}

func (o MultiValidate) Info() string {
	return o.check.Name
}

func (o *ValidationRequest) compareAnyValidator(checksData map[string]QueryResultInfo) (err error) {
	var firstName string
	var firstData QueryResults
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

func (o *ValidationRequest) compareValidator(checksData map[string]QueryResultInfo) error {
	return matchChecksData(checksData, o)
}

func (o *ValidateCheck) logBuildCheckNotPossible(err error) {
	Log.Info("Can't build check '%v' because of '%v'", o.Name, err)
}

func matchChecksData(checksData map[string]QueryResultInfo, req *ValidationRequest) (err error) {
	var firstName string
	var firstData QueryResults
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

func match(data1 QueryResults, data2 QueryResults, req *ValidationRequest) (err error) {
	var evalExpr *govaluate.EvaluableExpression
	if evalExpr, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	var results QueryResults
	if results, err = ComposeQueryResults("_", data1, data2); err == nil {
		err = validateData(results, evalExpr, true, "ib")
	}
	return
}
