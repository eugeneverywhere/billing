package handler

import (
	"unicode"
)

func ContainsSpaces(str string) bool {
	for _, v := range str {
		if unicode.IsSpace(v) {
			return true
		}
	}
	return false
}

func XORStrings(str1, str2 string) (result string) {
	for i := range str1 {
		result += string(str1[i] ^ str2[i%len(str2)])
	}
	return result
}
