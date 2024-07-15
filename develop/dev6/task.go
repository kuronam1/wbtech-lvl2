package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	fields    int
	delimiter string
	separated bool
)

func main() {
	flag.IntVar(&fields, "f", -1, "select column to show")
	flag.StringVar(&delimiter, "d", "\t", "choose another delimiter")
	flag.BoolVar(&separated, "s", false, "only strings with delimiter")

	flag.Parse()
	argument := flag.Args()

	if len(argument) != 1 {
		fmt.Fprint(os.Stderr, "cut signature:\n\tcat [-opts] filename")
		os.Exit(1)
	}

	file, err := os.Open(argument[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err while opening file: %v", err)
		os.Exit(1)
	}

	parseData(file)
}

func parseData(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if fields <= 0 {
			if separated && !strings.Contains(line, delimiter) {
				continue
			}
			fmt.Println(line)
		} else {
			separatedLine := strings.SplitAfter(line, delimiter)

			if len(separatedLine) < fields {
				if !separated {
					fmt.Println("")
				}
			} else {
				if separated && !strings.Contains(separatedLine[fields-1], delimiter) {
					continue
				}
				fmt.Println(separatedLine[fields-1])
			}
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "err while scanning: %v", err)
		os.Exit(1)
	}
}
