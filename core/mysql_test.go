package core

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
)

func TestMySqlService(t *testing.T) {
	service := MySqlService{mysql: &MySql{Name: "mysql", AccessKey: "mysql", Host: "localhost", Port: 3306},
		accessFinder:          AccessFinder()}
	err := service.Init()
	var check Check
	if check, err = service.NewСheck(&ValidationRequest{Query: "Select 1 as C1", EvalExpr: "C1 >= 1"}); err == nil {
		if data, err :=check.Query(); err == nil {
			println(data.String())
		}
		validateErr := check.Validate()
		AssertEqual(t, validateErr, nil, ErrorMessageBuilder)
	}

	if err != nil {
		println(err.Error())
	}
}
