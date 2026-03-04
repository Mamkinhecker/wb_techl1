package main

import "strings"

func isE(s string) bool {
	mp := make(map[rune]bool)
	s = strings.ToLower(s)

	data := []rune(s)
	for _, i := range data {
		if mp[i] {
			return false
		}

		mp[i] = true
	}
	return true
}
