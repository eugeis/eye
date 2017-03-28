package core

import (
	"strings"
	"fmt"
	"bufio"
	"os"
	"errors"
)

type Access struct {
	Key      string
	User     string
	Password string
}

type AccessFinder interface {
	FindAccess(key string) (Access, error)
}

type Security struct {
	Access []Access
}

func (o *Security) FindAccess(key string) (ret Access, err error) {
	for _, access := range o.Access {
		if strings.EqualFold(key, access.Key) {
			ret = access
			err = nil
			break
		}
	}
	if ret.Key == "" {
		err = errors.New(fmt.Sprintf("No access data found for '%v'", key))
	}
	return
}

func read() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
