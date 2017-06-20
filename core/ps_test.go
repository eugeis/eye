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
	var data QueryResults
	fmt.Printf("%v", time.Now())
	if check, err = ps.NewСheck(&ValidationRequest{}); err == nil {
		if data, err = check.Query(); err == nil {
			println(fmt.Sprintf("%v - %v", time.Now(), data.String()))
		}
	}

	if err != nil {
		println(err.Error())
	}

}
