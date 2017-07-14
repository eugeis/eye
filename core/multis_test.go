package core

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/eugeis/gee/as"
	"github.com/eugeis/eye/integ"
)

func TestMultiService(t *testing.T) {
	service1 := &MySqlService{mysql: &MySql{Name: "mysql", AccessKey: "mysql", Host: "localhost", Port: 3306},
		accessFinder:                AccessFinder()}
	if err := service1.Init(); err != nil {
		AssertEqual(t, err, nil, ErrorMessageBuilder)
		return
	}

	serviceFactory := NewFactory()
	eye := &Eye{config: &Config{}, accessFinder: &as.Security{}, serviceFactory: serviceFactory, liveChecks: integ.NewCache()}
	serviceFactory.add("mysql1", service1)
	serviceFactory.add("mysql2", service1)

	if check, err := eye.buildCompareCheck("test", []string{"mysql1", "mysql2"}, false,
		&ValidationRequest{Query: "Select 1 as C1", EvalExpr: "C1_1 == C1_2"}); err == nil {
		validateErr := check.Validate()
		AssertEqual(t, validateErr, nil, ErrorMessageBuilder)
	} else {
		AssertEqual(t, err, nil, ErrorMessageBuilder)
	}
}
