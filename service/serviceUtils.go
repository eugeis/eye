package service

import (
	"regexp"
	"time"
	"context"
	"rest/integ"
)

var log = integ.New("service: ")

func Match(pattern *regexp.Regexp, query func() ([]byte, error)) (ok bool, err error) {
	data, err := query()
	if err != nil {
		log.Info("Service failed because of '%s'", err)
	} else {
		log.Debug("data '%s'", data)
		ok = pattern.Match(data)
		if err != nil {
			log.Info("error '%s' for pattern '%s'\n", err, pattern)
		} else if !ok {
			log.Debug("no match for pattern '%s'\n", pattern)
		} else {
			log.Debug("match for pattern '%s'\n", pattern)
		}

	}
	return
}

func TimeoutContext(timeout time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), timeout)
	return c
}
