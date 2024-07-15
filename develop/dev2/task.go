package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const text = "qwe\\4\\5"

var invalidStringErr = fmt.Errorf("invalid string")

func unpackString(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	runes := []rune(s)

	var builder strings.Builder
	for i := 0; i < len(runes); i++ {
		current := runes[i]

		if i == len(runes)-1 {
			previous := runes[i-1]
			if unicode.IsDigit(current) && !unicode.IsLetter(previous) {
				return "", invalidStringErr
			} else {
				builder.Grow(1)
				builder.WriteRune(current)
			}
			break
		}

		next := runes[i+1]

		switch {
		case unicode.IsLetter(current) || unicode.IsPunct(current):
			if current == '\\' {
				i++
				current = runes[i]
				next = runes[i+1]
			}

			shift, err := unpackOne(next, current, &builder)
			if err != nil {
				return "", err
			}
			i += shift
		case unicode.IsDigit(current):
			return "", invalidStringErr
		}

	}

	return builder.String(), nil
}

func unpackOne(next, current rune, builder *strings.Builder) (int, error) {
	switch {
	case unicode.IsDigit(next):

		amount, err := strconv.Atoi(string(next))
		if err != nil {
			return 0, err
		}

		builder.Grow(amount)
		if err = unpackWithAmount(current, amount, builder); err != nil {
			return 0, err
		}

		return 1, nil
	default:
		builder.Grow(1)
		builder.WriteString(string(current))
	}

	return 0, nil
}

func unpackWithAmount(r rune, amount int, builder *strings.Builder) error {
	for i := 0; i < amount; i++ {
		_, err := builder.WriteRune(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	s, err := unpackString(text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}
