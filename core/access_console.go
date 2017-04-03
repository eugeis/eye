package core

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"github.com/howeyc/gopass"
)


func BuildAccessFinderFromConsole(config *Config) (ret AccessFinder, err error) {
	security := &Security{}
	ret = security
	security.Access = extractAccessKeys(config)
	fillAccessDataFromConsole(security)
	return
}

func fillAccessDataFromConsole(security *Security) (err error) {
	reader := bufio.NewReader(os.Stdin)
	var text string
	var pw []byte
	for key, item := range security.Access {
		fmt.Printf("Enter access data for '%v'\n", key)

		fmt.Print("User: ")
		text, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		item.User = strings.TrimSpace(text)

		fmt.Print("Password: ")
		pw, err = gopass.GetPasswdMasked()
		if err != nil {
			break
		}
		item.Password = string(pw)
		security.Access[key] = item
	}
	return
}
