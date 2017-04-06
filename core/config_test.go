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
	config, err := LoadConfig(
		"d:/views/diagnosis/config/eye.yml",
		"d:/views/diagnosis/config/eye.properties")
	if err != nil {
		println(err.Error())
	}
	println(config)
}
