package core

import (
	"testing"
)

func TestFillAccessData(t *testing.T) {
	security := &Security{}
	fillAccessData(security, "Security.yml")
	println(security)
}

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("d:/views/diagnosis/config/mega/serv1/eye.yml", "")
	if err != nil {
		println(err.Error())
	}
	println(config)
}
