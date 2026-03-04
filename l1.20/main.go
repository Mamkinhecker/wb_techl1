package main

func reverse(r []rune, left, right int) {
	for left < right {
		r[left], r[right] = r[right], r[left]
		left++
		right--
	}
}

func reverseWords(s string) string {
	runes := []rune(s)

	reverse(runes, 0, len(runes)-1)

	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || runes[i] == ' ' {
			reverse(runes, start, i-1)
			start = i + 1
		}
	}

	return string(runes)
}
