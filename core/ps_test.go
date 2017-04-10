package core

import (
	"testing"
)

func TestProcessService(t *testing.T) {
	ps := PsService{ps: &Ps{}}
	err := ps.Init()
	var check Check
	var data QueryResult
	if check, err = ps.NewСheck(&QueryRequest{}); err == nil {
		if data, err = check.Query(); err == nil {
			println(string(data))
		}
	}

	if err != nil {
		println(err.Error())
	}

}
