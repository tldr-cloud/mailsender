package main

import (
	"fmt"
	mailsender "tldr.cloud/mailsender"
)

func main() {
	if err := mailsender.SendWelcomeMail("viacheslav@kovalevskyi.com"); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("done")
}
