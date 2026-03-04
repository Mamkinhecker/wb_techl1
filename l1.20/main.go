package main

func ReverseStrings(s string) string {
	var result string
	for i := 0; i < len(s); i++ {
		j := i
		for ; j != len(s) && s[j] != ' '; j++ {
		}
		f := j + 1
		for ; j >= i; j-- {
			result += s[j]
		}

	}

	return result
}
