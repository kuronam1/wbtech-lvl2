package main

import (
	"fmt"
	"slices"
	"strings"
)

var test = []string{"пятак", "пятка", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

// AnagramsSearch ... написать тесты
func AnagramsSearch(array []string) map[string][]string {
	dictionary := make(map[string][]string)

	for _, value := range array {
		if key, found := findAnagram(value, dictionary); found {
			if !slices.Contains(dictionary[key], value) {
				dictionary[key] = append(dictionary[key], strings.ToLower(value))
			}
		} else {
			dictionary[value] = append(dictionary[value], strings.ToLower(value))
		}
	}

	sortAndSerialize(dictionary)

	return dictionary
}

func findAnagram(word string, dictionary map[string][]string) (string, bool) {
	sortedWord := sortString(word)

	for key := range dictionary {
		if len(word) == len(key) {
			sortedKey := sortString(key)
			if sortedWord == sortedKey {
				return key, true
			}
		}
	}

	return "", false
}

func sortAndSerialize(dictionary map[string][]string) {
	for key, variety := range dictionary {
		if len(variety) == 1 {
			delete(dictionary, key)
		} else {
			slices.Sort(variety)
		}
	}
}

func sortString(word string) string {
	stringInRunes := []rune(word)
	slices.Sort(stringInRunes)
	return string(stringInRunes)
}

func main() {
	fmt.Println(test)

	result := AnagramsSearch(test)
	for key, variety := range result {
		fmt.Printf("key: %s, variety: %v\n", key, variety)
	}
}
