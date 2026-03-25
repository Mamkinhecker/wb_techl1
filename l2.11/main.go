package main

import (
	"sort"
	"strings"
)

type pair struct {
	First  string
	Second string
}

func findAnagrams(words []string) map[string][]string {
	groups := make(map[string][]string)
	firstWord := make(map[string]string)

	for _, w := range words {
		lower := strings.ToLower(w)
		runes := []rune(lower)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)

		groups[key] = append(groups[key], w)
		if len(groups[key]) == 1 {
			firstWord[key] = w
		}
	}

	result := make(map[string][]string)
	for key, group := range groups {
		if len(group) >= 2 {
			sort.Strings(group)
			result[firstWord[key]] = group
		}
	}
	return result
}
