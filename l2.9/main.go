package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var builder strings.Builder
	runes := []rune(s)
	n := len(runes)
	i := 0
	for i < n {
		if runes[i] == '\\' {
			i++
			if i == n {
				return "", errors.New("wrong string!")
			}

			curr := runes[i]
			i++
			count, next, err := readCount(runes, i)
			if err != nil {
				return "", err
			}

			if count != 0 {
				builder.WriteString(strings.Repeat(string(curr), count))
			}

			i = next
			continue
		}

		if unicode.IsDigit(runes[i]) {
			return "", errors.New("wrong string!")
		}

		curr := runes[i]
		i++
		count, next, err := readCount(runes, i)
		if err != nil {
			return "", err
		}
		if count != 0 {
			builder.WriteString(strings.Repeat(string(curr), count))
		}

		i = next
	}
	return builder.String(), nil
}

func readCount(runes []rune, start int) (int, int, error) {
	if start >= len(runes) {
		return 1, start, nil
	}
	if !unicode.IsDigit(runes[start]) {
		return 1, start, nil
	}

	j := start
	for j < len(runes) && unicode.IsDigit(runes[j]) {
		j++
	}
	numStr := string(runes[start:j])
	count, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, 0, errors.New("invalid number")
	}
	return count, j, nil
}
