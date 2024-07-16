package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type cmd struct {
	command string
	args    []string
	stdin   chan string
	stdout  chan string
	stderr  chan error
	ctx     context.Context
}

func main() {
	for {
		fmt.Print("shell> ")
		inputLine, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while reading input: %v", err)
			os.Exit(1)
		}

		inputLine = strings.TrimSpace(inputLine)

		if inputLine == "\\q" {
			break
		}

		commands := strings.Split(inputLine, "|")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cmds := make([]cmd, 0, 2)
		errChan := make(chan error)
		for _, command := range commands {
			commandParts := strings.Fields(command)

			if len(commandParts) == 1 {
				cmds = append(cmds, cmd{
					command: commandParts[0],
					args:    nil,
					stdin:   nil,
					stdout:  nil,
					stderr:  errChan,
					ctx:     ctx,
				})
			} else {
				cmds = append(cmds, cmd{
					command: commandParts[0],
					args:    commandParts[1:],
					stdin:   nil,
					stdout:  nil,
					stderr:  errChan,
					ctx:     ctx,
				})
			}
		}

		var lastPipe chan string
		cmds[0].stdin, lastPipe = make(chan string), cmds[0].stdout

		for i := 1; i < len(cmds); i++ {
			cmds[i].stdin = lastPipe
			defer func() {

			}()
			lastPipe = cmds[i].stdout
		}
	}
}

func (c *cmd) Execute() error {
	return nil
}
