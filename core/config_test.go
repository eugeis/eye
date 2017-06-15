package core

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig([]string{
		"d:/views/diagnosis/config/eye.yml",
		"d:/views/diagnosis/config/eye.properties"}, []string{}, "d:/TC_CACHE/diagnosis")
	if err != nil {
		println(err.Error())
	}
	println(config)
}
