package conf

import (
	"reflect"
	"path/filepath"
	"bytes"
	"os"
	"strings"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/BurntSushi/toml"
	"errors"
	"time"
	"text/template"
)

func ReadConfig(config interface{}, file string, params map[string]string) (err error) {
	var data []byte
	if data, err = bindToProperties(file, params); err == nil {
		switch {
		case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
			err = yaml.Unmarshal(data, config)
		case strings.HasSuffix(file, ".toml"):
			err = toml.Unmarshal(data, config)
		case strings.HasSuffix(file, ".json"):
			err = json.Unmarshal(data, config)
		default:
			if toml.Unmarshal(data, config) != nil {
				if json.Unmarshal(data, config) != nil {
					if yaml.Unmarshal(data, config) != nil {
						err = errors.New("failed to decode config")
					}
				}
			}
		}
	}
	return
}

//from Hugo
func dfault(dflt interface{}, given ...interface{}) (interface{}, error) {
	// given is variadic because the following construct will not pass a piped
	// argument when the key is missing:  {{ index . "key" | default "foo" }}
	// The Go template will complain that we got 1 argument when we expectd 2.

	if len(given) == 0 {
		return dflt, nil
	}
	if len(given) != 1 {
		return nil, fmt.Errorf("wrong number of args for default: want 2 got %d", len(given)+1)
	}

	g := reflect.ValueOf(given[0])
	if !g.IsValid() {
		return dflt, nil
	}

	set := false

	switch g.Kind() {
	case reflect.Bool:
		set = true
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		set = g.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		set = g.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		set = g.Uint() != 0
	case reflect.Float32, reflect.Float64:
		set = g.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		set = g.Complex() != 0
	case reflect.Struct:
		switch actual := given[0].(type) {
		case time.Time:
			set = !actual.IsZero()
		default:
			set = true
		}
	default:
		set = !g.IsNil()
	}

	if set {
		return given[0], nil
	}

	return dflt, nil
}

func bindToProperties(file string, params map[string]string) (data []byte, err error) {
	tmpl := template.New(filepath.Base(file)).Funcs(template.FuncMap{
		"default": dfault,
	})
	if tmpl, err = tmpl.ParseFiles(file); err == nil {
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, params)
		data = buffer.Bytes()
	}
	return
}

func Env() (ret map[string]string) {
	ret = make(map[string]string)
	SplitParamsIntoMap(os.Environ(), ret)
	return
}

func SplitParamsIntoMap(params []string, fill map[string]string) {
	for _, e := range params {
		pair := strings.Split(e, "=")
		fill[pair[0]] = pair[1]
	}
}
