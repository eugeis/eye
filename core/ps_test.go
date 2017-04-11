package core

import (
	"testing"
	"time"
	"fmt"
)

func TestProcessService(t *testing.T) {
	ps := PsService{Ps: &Ps{}}
	err := ps.Init()
	var check Check
	var data QueryResult
	if check, err = ps.NewСheck(&QueryRequest{}); err == nil {
		for i := 1; i <= 10; i++ {
			if data, err = check.Query(); err == nil {
				println(fmt.Sprintf("%v - %v", time.Now(), string(data)))
			}
		}

	}

	if err != nil {
		println(err.Error())
	}

}
