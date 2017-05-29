package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"errors"
	"context"
	"regexp"
	"io"
)

type Service interface {
	Name() string

	Init() error
	Close()

	Ping() error

	NewСheck(req *QueryRequest) (Check, error)

	NewExporter(req *ExportRequest) (Exporter, error)
}

type Check interface {
	Info() string
	Query() (QueryResult, error)
	Validate() error
}

type Exporter interface {
	Info() string
	Export(params ...string) error
}

type Factory interface {
	Find(name string) (Service, error)
	Close()
}

type QueryRequest struct {
	Query string
	Expr  string
	Not   bool
}

type ExportRequest struct {
	Query     string
	Converter func(interface{}) string
	Out       func(interface{}) io.Writer
}

type CompareRequest struct {
	QueryRequest *QueryRequest
	Tolerance    int
}

func (o *CompareRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v[tolerance(%v),not(%v)]", o.QueryRequest.ChecksKey(strictness, serviceNames),
		o.Tolerance, o.QueryRequest.Not)
}

func (o *QueryRequest) CheckKey(serviceName string) string {
	if o.Not {
		return fmt.Sprintf("%v.q(%v).e(%v).not()", serviceName, o.Query, o.Expr)
	} else {
		return fmt.Sprintf("%v.q(%v).e(%v)", serviceName, o.Query, o.Expr)
	}
}

func (o *QueryRequest) ChecksKey(strictness string, serviceNames []string) string {
	if o.Not {
		return fmt.Sprintf("%v(%v.q(%v).e(%v).not())", strictness, strings.Join(serviceNames, "-"), o.Query, o.Expr)
	} else {
		return fmt.Sprintf("%v(%v.q(%v).e(%v))", strictness, strings.Join(serviceNames, "-"), o.Query, o.Expr)
	}
}

func (o *ExportRequest) ExportKey(serviceName string) string {
	return fmt.Sprintf("%v.q(%v)", serviceName, o.Query)
}

type QueryResult []byte

func (o *QueryResult) String() (ret string, err error) {
	ret = string(*o)
	return
}

func (o *QueryResult) Int() (ret int, err error) {
	ret, err = strconv.Atoi(string(*o))
	return
}

func TimeoutContext(timeout time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), timeout)
	return c
}

type SimpleServiceFactory struct {
	services map[string]Service
}

func NewFactory() SimpleServiceFactory {
	return SimpleServiceFactory{services: make(map[string]Service)}
}

func (o *SimpleServiceFactory) Find(name string) (ret Service, err error) {
	nameLowCase := strings.ToLower(name)
	ret = o.services[nameLowCase]
	if ret == nil {
		err = errors.New(fmt.Sprintf("Can't find the service '%v'", name))
	}
	return
}

func (o *SimpleServiceFactory) Add(service Service) {
	nameLowCase := strings.ToLower(service.Name())
	o.services[nameLowCase] = service
}

func (o *SimpleServiceFactory) Close() {
	for _, item := range o.services {
		item.Close()
	}
	o.services = make(map[string]Service)
}

func compilePattern(pattern string) (ret *regexp.Regexp, err error) {
	if len(pattern) > 0 {
		ret, err = regexp.Compile(pattern)
	}
	return
}

func compilePatterns(pattern ...string) (ret []*regexp.Regexp, err error) {
	ret = make([]*regexp.Regexp, len(pattern))
	for _, p := range pattern {
		if len(p) > 0 {
			ret[0], err = regexp.Compile(p)
			if err != nil {
				break
			}
		}
	}
	return
}

func ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v(%v)", strictness, strings.Join(serviceNames, "-"))
}

func validate(check Check, pattern *regexp.Regexp, not bool) (err error) {
	var data QueryResult
	if data, err = check.Query(); err == nil {
		if pattern != nil {
			if match := pattern.Match(data); (!match && !not) || (match && not) {
				err = errors.New(fmt.Sprintf("No match for %v", check.Info()))
			}
		} else if match := data != nil; (!match && !not) || (match && not) {
			err = errors.New(fmt.Sprintf("No match, empty result for %v", check.Info()))
		}
	}
	return
}
