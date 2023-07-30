package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func askForConfirmation(s string) bool {
  reader := bufio.NewReader(os.Stdin)

  	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
