package conf

import (
	"regexp"
	"strings"
	"bytes"
	"errors"
	"os"
)

type PropsReader struct {
	pattern *regexp.Regexp
}

func New() PropsReader {
	pattern, _ := regexp.Compile("[#].*\\n|\\s+\\n|\\S+[=]|.*\n")
	return PropsReader{pattern: pattern}
}

var props = New()

func ParseProperties(data bytes.Buffer) (map[string]string, error) {
	return ParsePropertiesInto(data, make(map[string]string))
}
func ParsePropertiesInto(data bytes.Buffer, fill map[string]string) (ret map[string]string, err error) {
	ret = fill
	str := data.String()
	if !strings.HasSuffix(str, "\n") {
		str = str + "\n"
	}
	s2 := props.pattern.FindAllString(str, -1)

	for i := 0; i < len(s2); {
		if strings.HasPrefix(s2[i], "#") {
			i++
		} else if strings.HasSuffix(s2[i], "=") {
			key := s2[i][0: len(s2[i])-1]
			i++
			if strings.HasSuffix(s2[i], "\n") {
				val := s2[i][0: len(s2[i])-1]
				if strings.HasSuffix(val, "\r") {
					val = val[0: len(val)-1]
				}
				i++
				fill[key] = val
			}
		} else if strings.Index(" \t\r\n", s2[i][0:1]) > -1 {
			i++
		} else {
			err = errors.New("Unable to process line in cfg file containing " + s2[i])
			break
		}
	}
	return
}

func Env() map[string]string {
	return EnvInto(make(map[string]string))
}

func EnvInto(fill map[string]string) (ret map[string]string) {
	ret = fill
	SplitPropertiesIntoMap(os.Environ(), ret)
	return
}

func SplitPropertiesIntoMap(params []string, fill map[string]string) {
	for _, e := range params {
		pair := strings.Split(e, "=")
		fill[pair[0]] = pair[1]
	}
}
