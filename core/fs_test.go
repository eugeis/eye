package core

import (
	"testing"
	"time"
	"fmt"
)

func TestFileSystemService(t *testing.T) {
	service := FsService{Fs: &Fs{File:"D:/views"}}
	err := service.Init()
	var check Check
	var data QueryResult
	if check, err = service.NewСheck(&ValidationRequest{}); err == nil {
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
