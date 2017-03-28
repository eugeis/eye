package core

import (
	"github.com/jinzhu/configor"
	"strings"
	"fmt"
	"bufio"
	"os"
	"errors"
)

type AccessFinder interface {
	FindAccess(key string) (Access, error)
}

type Security struct {
	Access []Access
}

func LoadSecurity(fileNames []string) (ret *Security, err error) {
	ret = &Security{}
	err = configor.Load(ret, fileNames...)

	//ignore, https://github.com/jinzhu/configor/issues/6
	if err != nil && strings.EqualFold(err.Error(), "invalid config, should be struct") {
		err = nil
	}
	return
}

func (c *Security) FindAccess(key string) (ret Access, err error) {
	err = errors.New(fmt.Sprintf("No access data found for '%s'", key))
	for _, access := range c.Access {
		if strings.EqualFold(key, access.Key) {
			ret = access
			err = nil
			break
		}
	}
	return
}

func read() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
