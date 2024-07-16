package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const ignoreCasePrefix = "(?i)"

var (
	after      int  //+
	before     int  //+
	ctx        int  //+
	count      bool //+
	ignoreCase bool //+
	invert     bool //+
	fixed      bool //+
	lineNum    bool //+
	counter    int
)

func main() {
	flag.IntVar(&after, "A", -1, "print N strings after match")
	flag.IntVar(&before, "B", -1, "print N strings before match")
	flag.IntVar(&ctx, "C", -1, "print N strings before and after match")
	flag.BoolVar(&count, "c", false, "count matches")
	flag.BoolVar(&ignoreCase, "i", false, "ignore letter case")
	flag.BoolVar(&invert, "v", false, "print all that not matches expression")
	flag.BoolVar(&fixed, "F", false, "find exact match with string")
	flag.BoolVar(&lineNum, "n", false, "print lines with match")

	flag.Parse()

	if ctx > 0 && (after > 0 || before > 0) {
		fmt.Fprintln(os.Stderr, "cannot use flag -C with -A or -B")
		os.Exit(1)

	}

	arguments := flag.Args()
	if len(arguments) != 2 {
		fmt.Fprintln(os.Stderr, "usage:\n\t./task [options] expression filename")
		os.Exit(1)
	}

	file, err := getFile(arguments[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	fileData := make([]string, 0, 10)
	for scanner.Scan() {
		line := scanner.Text()

		fileData = append(fileData, line)
	}

	expressionString := strings.TrimSpace(arguments[0])
	if ignoreCase && !strings.Contains(expressionString, ignoreCasePrefix) {
		expressionString = ignoreCasePrefix + expressionString
	}

	expression := regexp.MustCompile(expressionString)

	processingData(fileData, expression)
}

func processingData(fileData []string, expression *regexp.Regexp) {
	for idx, line := range fileData {
		if expression.MatchString(line) && !invert {
			processLine(idx, fileData)
		}

		if !expression.MatchString(line) && invert {
			processLine(idx, fileData)
		}
	}

	if count && counter != 0 {
		fmt.Fprintf(os.Stdout, "lines counted: %v", counter)
	} else if counter == 0 {
		fmt.Fprintln(os.Stdout, "no lines are matching pattern or string")
	}
}

func processLine(idx int, fileData []string) {
	counter++

	if ctx > 0 && (idx-ctx) > -1 && (idx+ctx) < len(fileData) {
		printCtx(idx, fileData)
	} else {
		if before > 0 && (idx-before) > -1 {
			printBefore(idx, fileData)
		}

		printLine(idx, fileData[idx])

		if after > 0 && (idx+after) < len(fileData) {
			printAfter(idx, fileData)
		}
	}
}

func printLine(idx int, line string) {
	if lineNum {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%d:%s", idx+1, line))
	} else {
		fmt.Fprintln(os.Stdout, line)
	}
}

func printCtx(idx int, fileData []string) {
	printBefore(idx, fileData)
	printLine(idx, fileData[idx])
	printAfter(idx, fileData)
}

func printBefore(idx int, fileData []string) {
	for i := idx - before; i < idx; i++ {
		printLine(i, fileData[i])
	}
}

func printAfter(idx int, fileData []string) {
	for i := idx + 1; i <= idx+after; i++ {
		printLine(i, fileData[i])
	}
}

func getFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("file with path %v does not exists", filepath))
		} else {
			return nil, err
		}
	}

	fileStats, err := file.Stat()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("err while opening file: %v", err))
	}

	if fileStats.IsDir() {
		return nil, errors.New(fmt.Sprintf("%v is a directory", filepath))
	}

	return file, nil
}
