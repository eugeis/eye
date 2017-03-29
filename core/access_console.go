package core

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"github.com/howeyc/gopass"
)

func fillAccessDataFromConsole(security *Security) {
	reader := bufio.NewReader(os.Stdin)
	for key, item := range security.Access {
		fmt.Printf("Enter access data for '%v'\n", key)
		fmt.Print("User: ")
		text, err := reader.ReadString('\n')
		if err == nil {
			item.User = strings.TrimSpace(text)
		}
		fmt.Print("Password: ")
		pw, err := gopass.GetPasswdMasked()
		if err == nil {
			item.Password = pw
		}
		security.Access[key] = item
	}
}
