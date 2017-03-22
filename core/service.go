package core

import (
	"fmt"
	"strings"
)

type Service interface {
	Name() string
	Kind() string

	Init() error
	Close()

	Ping() error

	NewСheck(req *QueryRequest) (Check, error)
}

type Check interface {
	Info() string
	Query() ([]byte, error)
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
	Tolerance     float64
	Not bool
}

func (o QueryRequest) CheckKey(serviceName string) string {
	return fmt.Sprintf("%s.q(%s).e(%s)", serviceName, o.Query, o.Expr)
}

func (o QueryRequest) ChecksKey(serviceNames []string) string {
	return fmt.Sprintf("%s.q(%s).e(%s)", strings.Join(serviceNames, "_"), o.Query, o.Expr)
}
