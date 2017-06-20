package core

import (
	"testing"
	"time"
	"fmt"
)

func TestFileSystemService(t *testing.T) {
	service := FsService{Fs: &Fs{File: "D:/TC_CACHE"}}
	err := service.Init()
	var check Check
	var data QueryResults
	if check, err = service.NewСheck(&ValidationRequest{}); err == nil {
		if data, err = check.Query(); err == nil {
			println(fmt.Sprintf("%v - %v", time.Now(), data.String()))
		}
	}

	if err != nil {
		println(err.Error())
	}

}
