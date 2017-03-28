package core

import (
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
