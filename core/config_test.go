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
	config, _ := Load([]string{"d:/views/share/t/diagnosis/config/mega/serv1/eye.yml"})
	println(config)
}
