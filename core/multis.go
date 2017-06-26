package core

import "gopkg.in/Knetic/govaluate.v2"

type MultiCheck struct {
	info        string
	queries     []Query
	eval        *govaluate.EvaluableExpression
	all         bool
	onlyRunning bool
}

func (o *MultiCheck) Validate() (err error) {
	var data QueryResults
	if data, err = o.checksData(); err == nil {
		err = validateData(data, o.eval, o.all, o.info)
	}
	return
}

func (o *MultiCheck) Query() (data QueryResults, err error) {
	return
}

func (o *MultiCheck) Info() string {
	return o.info
}

func (o *MultiCheck) checksData() (ret QueryResults, err error) {
	querysData := make([]QueryResults, 0)
	for _, check := range o.queries {
		data, queryErr := check.Query()
		if queryErr == nil {
			querysData = append(querysData, data)
		} else if o.onlyRunning {
			err = queryErr
			break
		}
	}
	if err == nil {
		ret, err = ComposeQueryResults("_", querysData)
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

func (o *MultiPing) Query() (data QueryResults, err error) {
	return
}

func (o *MultiPing) Info() string {
	return o.check.Name
}

type MultiValidate struct {
	check     *ValidateCheck
	validator func([]string, *ValidationRequest) error
}

func (o *MultiValidate) Validate() (err error) {
	return o.validator(o.check.Services, o.check.Request)
}

func (o *MultiValidate) Query() (data QueryResults, err error) {
	return
}

func (o *MultiValidate) Info() string {
	return o.check.Name
}

func (o *ValidateCheck) logBuildCheckNotPossible(err error) {
	Log.Info("Can't build check '%v' because of '%v'", o.Name, err)
}
