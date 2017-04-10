package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"errors"
	"context"
	"regexp"
)

type Service interface {
	Name() string

	Init() error
	Close()

	Ping() error

	NewСheck(req *QueryRequest) (Check, error)
}

type Check interface {
	Info() string
	Query() (QueryResult, error)
	Validate() error
}

type Factory interface {
	Find(name string) (Service, error)
	Close()
}

type QueryRequest struct {
	Query string
	Expr  string
}

type CompareRequest struct {
	QueryRequest *QueryRequest
	Tolerance    int
	Not          bool
}

func (o *CompareRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v[tolerance(%v),not(%v)]", o.QueryRequest.ChecksKey(strictness, serviceNames),
		o.Tolerance, o.Not)
}

func (o *QueryRequest) CheckKey(serviceName string) string {
	return fmt.Sprintf("%v.q(%v).e(%v)", serviceName, o.Query, o.Expr)
}

func (o *QueryRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v(%v.q(%v).e(%v))", strictness, strings.Join(serviceNames, "-"), o.Query, o.Expr)
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

func compileRegexp(req *QueryRequest) (ret *regexp.Regexp, err error) {
	if len(req.Expr) > 0 {
		ret, err = regexp.Compile(req.Expr)
	}
	return
}
