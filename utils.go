package main

import (
	"fmt"
	"regexp"
)

func matchPattern(pattern string) func(name string) bool {
	reg := regexp.MustCompile(pattern)
	return func(name string) bool {
		return reg.MatchString(name)
	}
}

func filterStrings(list []string, pattern string) []string {
	cleaned := []string{}
	match := matchPattern(pattern)
	for _, elm := range list {
		if match(elm) {
			cleaned = append(cleaned, elm)
		}
	}
	return cleaned
}

func printStrings(list []string) {
	for _, elm := range list {
		fmt.Println(elm)
	}
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}
