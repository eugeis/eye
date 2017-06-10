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

	NewСheck(req *ValidationRequest) (Check, error)

	NewExporter(req *ExportRequest) (Exporter, error)
}

type Check interface {
	Info() string
	Query() (QueryResult, error)
	Validate() error
}

type Exporter interface {
	Info() string
	Export(params map[string]string) error
}

type Factory interface {
	Find(name string) (Service, error)
	Close()
}

type ValidationRequest struct {
	Query     string
	Expr      string
	Operator  string
	Tolerance int
	Not       bool
}

type ExportRequest struct {
	Query     string
	Converter func(map[string]interface{}) []byte
	Out       func(params map[string]string) (io.WriteCloser, error)
}

type CompareRequest struct {
	ValidationRequest *ValidationRequest
	Tolerance         int
}

func (o *CompareRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v[tolerance(%v)]", o.ValidationRequest.ChecksKey(strictness, serviceNames),
		o.Tolerance)
}

func (o *ValidationRequest) CheckKey(serviceName string) string {
	return fmt.Sprintf("%v.q(%v).e(%v).op(%v).tolerance(%v).not(%v)",
		serviceName, o.Query, o.Expr, o.Operator, o.Tolerance, o.Not)
}

func (o *ValidationRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v(%v.q(%v).e(%v).op(%v).tolerance(%v).not(%v))", strictness,
		strings.Join(serviceNames, "-"), o.Query, o.Expr, o.Operator, o.Tolerance, o.Not)
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

func prepareQuery(query string, params map[string]string) (ret string) {
	ret = query
	if params != nil && len(params) > 0 {
		for k, v := range params {
			placeHolder := fmt.Sprintf("@@%v@@", strings.ToUpper(k))
			ret = strings.Replace(ret, placeHolder, v, -1)
		}
	}
	return
}
