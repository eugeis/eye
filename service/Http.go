package service

import (
	"regexp"
	"time"
	"fmt"
	"io/ioutil"
	"rest/digest"
	"net/http"
	"errors"
)

type Http struct {
	ServiceName string
	Url         string

	User     string
	Password string

	PingTimeoutMillis  int
	QueryTimeoutMillis int

	client *http.Client
	req    *digest.Request

	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (s *Http) Kind() string {
	return "Http"
}

func (s *Http) Name() string {
	return s.ServiceName
}

type httpCheck struct {
	info    string
	req     *digest.Request
	pattern *regexp.Regexp
	service *Http
}

func (o httpCheck) Check() (ret bool, err error) {
	err = o.service.Init()
	if err == nil {
		ret, err = Match(o.info, o.pattern, func() (data []byte, err error) {
			resp, err := o.req.Execute(o.service.client)
			if err == nil {
				data, err = ioutil.ReadAll(resp.Body)
				log.Debug("%s", data)
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					err = errors.New(fmt.Sprintf("Status %d", resp.StatusCode))
				}
			}
			return
		})
	}
	return

}

func (s *Http) NewСheck(query string, expr string) (ret Check, err error) {
	pattern, err := regexp.Compile(expr)
	if err == nil {
		req := digest.NewRequest(s.User, s.Password, "GET", s.Url+query, "")
		ret = httpCheck{
			info:    fmt.Sprintf("q: %s, e: %s", query, expr),
			req:     &req,
			pattern: pattern, service: s }
	}
	return
}

func (s *Http) Init() error {
	var err error

	if s.req == nil {
		s.client = digest.NewClient(true, s.queryTimeout)
		req := digest.NewRequest(s.User, s.Password, "GET", s.Url, "")
		s.req = &req
	}
	return err
}

func (s *Http) Close() {
	if s.req != nil {
		s.req.Close()
	}
	s.req = nil
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
	resp, err := s.req.Execute(s.client)
	if err == nil {
		log.Debug(body(resp))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = errors.New(fmt.Sprintf("Status %s", resp.StatusCode))
		}
	}
	return err
}

func body(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	ret := fmt.Sprintf("%s", body)
	return ret
}
