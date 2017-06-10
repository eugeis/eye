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
	fmt.Printf("%v", time.Now())
	if check, err = ps.NewСheck(&ValidationRequest{}); err == nil {
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
