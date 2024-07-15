package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var inputString string
	for {
		fmt.Print("shell> ")
		_, err := fmt.Scan(&inputString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err while scanning: %v", err)
			os.Exit(1)
		}

		if inputString == "\\q" {
			break
		}

		commandPool := strings.Split(inputString, "|")
		for _, command := range commandPool {

		}
	}
}

func performCommand(command string) {

}
