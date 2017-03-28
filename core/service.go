package core

import (
	"fmt"
	"strconv"
	"strings"
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

func (o CompareRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%s[tolerance(%v),not(%v)]", o.QueryRequest.ChecksKey(strictness, serviceNames),
		o.Tolerance, o.Not)
}

func (o QueryRequest) CheckKey(serviceName string) string {
	return fmt.Sprintf("%s.q(%s).e(%s)", serviceName, o.Query, o.Expr)
}

func (o QueryRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%s(%s.q(%s).e(%s))", strictness, strings.Join(serviceNames, "-"), o.Query, o.Expr)
}

type QueryResult []byte

func (o QueryResult) String() (ret string, err error) {
	ret = string(o)
	return
}

func (o QueryResult) Int() (ret int, err error) {
	ret, err = strconv.Atoi(string(o))
	return
}
