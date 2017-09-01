package core

import (
	"github.com/eugeis/gee/as"
	"fmt"
	"testing"
)

func AccessFinder() as.AccessFinder {
	return &as.Security{data: map[string]as.Access{"mysql": as.Access{User: "root", Password: "***"}}}
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, messageBuilder func(a interface{}, b interface{}) string) {
	if a == b {
		return
	}
	var message string
	if messageBuilder == nil {
		message = DefaultMessageBuilder(a, b)
	} else {
		message = messageBuilder(a, b)
	}
	t.Fatal(message)
}

func DefaultMessageBuilder(a interface{}, b interface{}) string {
	return fmt.Sprintf("%v != %v", a, b)
}

func ErrorMessageBuilder(a interface{}, b interface{}) string {
	if aAsError, ok := a.(error); ok {
		return "Validation error occurred: " + aAsError.Error()
	} else {
		return DefaultMessageBuilder(a, b)
	}
}
