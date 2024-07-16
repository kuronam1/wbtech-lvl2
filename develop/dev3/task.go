package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type sorting interface {
	string | int
}

var (
	ignoreTales   bool
	checkSort     bool
	skipRepeat    bool
	reverse       bool
	numericSort   bool
	column        int
	delimiter     = " "
	argumentErr   = errors.New("need one file\nusage: ./task [options] filename")
	convertingErr = errors.New("err while converting number")
)

func main() {
	flag.IntVar(&column, "k", -1, "choose a column to sort")
	flag.BoolVar(&numericSort, "n", false, "sort by numeric value")
	flag.BoolVar(&reverse, "r", false, "sort by reverse order")
	flag.BoolVar(&skipRepeat, "u", false, "skip string repeats")
	flag.BoolVar(&checkSort, "c", false, "check if file is sorted")
	flag.BoolVar(&ignoreTales, "b", false, "ignore tale spaces")

	flag.Parse()

	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Fprintln(os.Stderr, argumentErr)
		os.Exit(1)
	}

	file, err := os.Open(arguments[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "err while opening file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if checkSort {

	}

	if column > 0 {
		columnSorting(scanner)
	} else {
		defaultSorting(scanner)
	}
}

func columnSorting(scanner *bufio.Scanner) {
	data := make([][]string, 0, 10)
	tmpData := make([][]string, 0, 10)

	for scanner.Scan() {
		lineSlice := strings.Split(scanner.Text(), delimiter)

		if len(lineSlice) < column {
			tmpData = append(tmpData, lineSlice)
			continue
		}

		if ignoreTales {
			lineSlice[column-1] = strings.TrimSpace(lineSlice[column-1])
		}

		if skipRepeat && containsElem(data, lineSlice) {
			continue
		}

		data = append(data, lineSlice)
	}

	if checkSort {
		for i := 0; i < len(data)-1; i++ {
			if reverse {
				if data[i][column-1] < data[i+1][column-1] {
					fmt.Fprintln(os.Stdout, "file is not sorted")
					return
				}
			} else {
				if data[i][column-1] > data[i+1][column-1] {
					fmt.Fprintln(os.Stdout, "file is not sorted")
					return
				}
			}
		}

		fmt.Fprintln(os.Stdout, "file is sorted")
		return
	}

	sort.Slice(data, func(i, j int) bool {
		if numericSort {
			num1 := mustAtoi(data[i][column-1])
			num2 := mustAtoi(data[j][column-1])

			if reverse {
				return num1 > num2
			} else {
				return num1 < num2
			}
		} else {
			element1 := data[i][column-1]

			element2 := data[j][column-1]

			if reverse {
				return element1 > element2
			} else {
				return element1 < element2
			}
		}
	})

	data = append(data, tmpData...)

	for _, line := range data {
		stringLine := strings.Join(line, delimiter)

		fmt.Fprintln(os.Stdout, stringLine)
	}
}

func containsElem(data [][]string, line []string) bool {
	for _, element := range data {
		if element[column-1] == line[column-1] {
			return true
		}
	}

	return false
}

func defaultSorting(scanner *bufio.Scanner) {
	if numericSort {
		data := make([]int, 0, 10)

		for scanner.Scan() {
			line := scanner.Text()

			if ignoreTales {
				line = strings.TrimSpace(line)
			}

			num := mustAtoi(line)

			if skipRepeat && slices.Contains(data, num) {
				continue
			}

			data = append(data, num)
		}

		if checkSort {
			if slices.IsSorted(data) {
				fmt.Fprintln(os.Stdout, "file is sorted")
			} else {
				fmt.Fprintln(os.Stdout, "file is not sorted")
			}

			return
		}

		sortAndOut(data)

	} else {
		data := make([]string, 0, 10)

		for scanner.Scan() {
			line := scanner.Text()

			if ignoreTales {
				line = strings.TrimSpace(line)
			}

			if skipRepeat && slices.Contains(data, line) {
				continue
			}

			data = append(data, line)
		}

		if checkSort {
			if slices.IsSorted(data) {
				fmt.Fprintln(os.Stdout, "file is sorted")
			} else {
				fmt.Fprintln(os.Stdout, "file is not sorted")
			}

			return
		}

		sortAndOut(data)
	}
}

func sortAndOut[S sorting](data []S) {
	sort.Slice(data, func(i, j int) bool {
		if reverse {
			return data[i] > data[j]
		} else {
			return data[i] < data[j]
		}
	})

	for _, number := range data {
		fmt.Fprintln(os.Stdout, number)
	}
}

func mustAtoi(ascii string) int {
	num, err := strconv.Atoi(ascii)
	if err != nil {
		panic(fmt.Errorf("%v: %v", convertingErr, err))
	}

	return num
}
