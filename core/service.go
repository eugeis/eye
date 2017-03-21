package core

import "fmt"

type Service interface {
	Name() string
	Kind() string

	Init() error
	Close()

	Ping() error

	NewСheck(req *QueryRequest) (Check, error)
}

type Check interface {
	Query() ([]byte, error)
	Validate() error
	Close()
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
	Operator     string
}

func (o QueryRequest) CommandKey(serviceName string) string {
	return fmt.Sprintf("%s.q(%s).e(%s)", serviceName, o.Query, o.Expr)
}
