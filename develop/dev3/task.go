package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
)

type compare interface {
	string | int
}

type line struct {
	l []string
}

type dataForColumnSort struct {
	data   []line
	column int
}

type dataForDefaultSort []string

var (
	skipRepeat  bool
	reverse     bool
	numericSort bool
	column      int
	delimiter   = " "
	argumentErr = errors.New("need one file\nusage: ./task [options] filename")
)

// Переделать все от 54 строчки
func main() {
	flag.IntVar(&column, "k", -1, "specify the column to sort")
	flag.BoolVar(&numericSort, "n", false, "sort in numeric order")
	flag.BoolVar(&reverse, "r", false, "sort in reverse order")
	flag.BoolVar(&skipRepeat, "u", false, "do not show repeats")

	flag.Parse()

	arguments := flag.Args()
	if len(arguments) != 1 {
		log.Fatal(argumentErr)
	}

	file, err := os.Open(arguments[0])
	if err != nil {
		log.Fatal("error while opening file: ", err)
	}
	defer file.Close()

	sortFile(file)
}

func sortFile(file *os.File) {
	scanner := bufio.NewScanner(file)

	if column != -1 {
		if numericSort {
			data := make([]int, 0, 10)
		} else {
			data := make([]string, 0, 10)
		}

		for scanner.Scan() {
			str := scanner.Text()

			if skipRepeat && slices.Contains(data, str) {
				continue
			}

			if numericSort {

			}

			data = append(data, str)
		}

	} else {
		lines := make([]string, 0, 10)
		for scanner.Scan() {
			str := scanner.Text()
			if skipRepeat && slices.Contains(lines, str) {
				continue
			}
			lines = append(lines, str)
		}

		if reverse {
			slices.Reverse(lines)
		}

		slices.Sort(lines)
	}
}

func readData[C compare](dataSlice []C, scanner *bufio.Scanner) {
	for scanner.Scan()
}

func (d *dataForColumnSort) Len() int {
	return len(d.data)
}

func (d *dataForColumnSort) Less(i, j int) bool {
	return d.data[i].l[d.column] < d.data[j].l[d.column]
}

func (d *dataForColumnSort) Swap(i, j int) {
	d.data[i], d.data[j] = d.data[j], d.data[i]
}

func (d *dataForColumnSort) Contains(ln string) bool {
	for _, value := range d.data {
		if strings.Join(value.l, delimiter) == ln {
			return true
		}
	}

	return false
}

func (d *dataForColumnSort) Sort() {
	if reverse {
		sort.Sort(sort.Reverse(d))
	} else {
		sort.Sort(d)
	}
}
