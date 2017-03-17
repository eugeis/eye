package service

import (
	"regexp"
	"time"
	"fmt"
	"io/ioutil"
	"rest/digest"
	"net/http"
	"github.com/pkg/errors"
)

type Http struct {
	ServiceName string
	Url         string

	User     string
	Password string

	PingTimeoutMillis  int
	QueryTimeoutMillis int

	req *digest.Request

	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (s *Http) Kind() string {
	return "Http"
}

func (s *Http) Name() string {
	return s.ServiceName
}

func (s *Http) Check(query string, pattern *regexp.Regexp) (ok bool, err error) {
	err = s.Init()
	if err != nil {
		return false, err
	}

	return Match(pattern, func() (data []byte, err error) {

		resp, err := s.req.ExecuteQuery(query)
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

func (s *Http) Init() error {
	var err error

	if s.req == nil {
		req := digest.NewRequest(s.User, s.Password, "GET", s.Url, "")

		resp, err := req.Execute()
		if err != nil {
			return err
		} else {
			log.Debug(body(resp))
		}
		defer resp.Body.Close()
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
	resp, err := s.req.Execute()
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
