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

func isValidLocale(s string) bool {
	// TODO: improve using the package "golang.org/x/text/language"
	if len(s) != 2 {
		return false
	}

	for _, c := range s {
		if c < 'a' || c > 'z' {
			return false
		}
	}

	return true
}
