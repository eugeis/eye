package core

import (
	"errors"
	"eye/digest"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Http struct {
	ServiceName string
	AccessKey   string
	Url         string

	PingRequest *QueryRequest

	PingTimeoutMillis  int
	QueryTimeoutMillis int

	client    *http.Client
	pingCheck *httpCheck

	pingTimeout  time.Duration
	queryTimeout time.Duration

	access func(string) Access
}

func (s *Http) Kind() string {
	return "Http"
}

func (s *Http) Name() string {
	return s.ServiceName
}

func (s *Http) Init() (err error) {
	if s.client == nil {
		s.client = digest.NewClient(true, s.queryTimeout)
		s.pingCheck, err = s.newСheck(s.PingRequest)
		if err != nil {
			s.Close()
		}
	}
	return
}

func (s *Http) Close() {
	s.client = nil
	s.pingCheck = nil
}

func (s *Http) Ping() error {
	err := s.Init()
	if err == nil {
		err = s.ping()
		if err != nil {
			log.Debug("'%s' can't be reached because of %s", s.Name(), err)
		}
	}
	return err
}

func (s *Http) ping() error {
	return s.pingCheck.Validate()
}

func body(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	ret := fmt.Sprintf("%s", body)
	return ret
}

func (s *Http) NewСheck(req *QueryRequest) (ret Check, err error) {
	return s.newСheck(req)
}

func (s *Http) newСheck(req *QueryRequest) (ret *httpCheck, err error) {
	access := s.access(s.AccessKey)
	if access.Key == s.AccessKey {
		var pattern *regexp.Regexp
		if len(req.Expr) > 0 {
			pattern, err = regexp.Compile(req.Expr)
			if err != nil {
				return
			}
		}

		dReq := digest.NewRequest(access.User, access.Password, "GET", s.Url+req.Query, "")
		ret = &httpCheck{
			info:    req.CheckKey(s.Name()),
			req:     &dReq,
			pattern: pattern, service: s}
	} else {
		err = errors.New(fmt.Sprintf("No access data found for the service %s", s.ServiceName))
	}
	return
}

//buildCheck

type httpCheck struct {
	info    string
	req     *digest.Request
	pattern *regexp.Regexp
	service *Http
}

func (o httpCheck) Info() string {
	return o.info
}

func (o httpCheck) Validate() (err error) {
	data, err := o.Query()

	if err == nil && o.pattern != nil {
		if !o.pattern.Match(data) {
			err = errors.New(fmt.Sprintf("No match for %s", o.info))
		}
	}
	return
}

func (o httpCheck) Query() (data QueryResult, err error) {
	err = o.service.Init()
	if err != nil {
		return
	}
	resp, err := o.req.Execute(o.service.client)
	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Debug("%s", data)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Status %d", resp.StatusCode))
	}
	return
}
