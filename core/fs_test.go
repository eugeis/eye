package core

import (
	"testing"
)

func TestFileSystemService(t *testing.T) {
	service := FsService{Fs: &Fs{File: "D:/TC_CACHE/CG-latest/server"}}
	err := service.Init()
	var check Check
	if check, err = service.NewСheck(NewValidationRequest("wildfly/standalone/deployments", "Name !~ '.*(\\.deploying|\\.failed)'")); err == nil {
		validateErr := check.Validate()
		AssertEqual(t, validateErr, nil, ErrorMessageBuilder)
	}

	if err != nil {
		println(err.Error())
	}

}
