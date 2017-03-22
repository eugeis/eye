package core

import "regexp"

type MultiCheck struct {
	info      string
	query     string
	pattern   *regexp.Regexp
	checks    []Check
	validator func(map[string][]byte) error
}

func (o MultiCheck) Validate() error {
	return o.validator(o.checksData())
}

func (o MultiCheck) Query() (data []byte, err error) {
	return
}

func (o MultiCheck) Info() string {
	return o.info
}

func (o MultiCheck) checksData() (checksData map[string][]byte) {
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
