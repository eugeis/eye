package core

import (
	"regexp"
	"github.com/pkg/errors"
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
				} else {
					checksData[check.Info()] = QueryResultInfo{result: matches[0]}
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
