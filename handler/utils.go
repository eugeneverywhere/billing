package handler

import (
	"strings"
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

func XORStrings(str1, str2 string) string {
	var short, long string
	if len(str1) > len(str2) {
		short = str2
		long = str1
	} else {
		short = str1
		long = str2
	}
	result := make([]byte, len(long))
	short += strings.Repeat(" ", len(long)-len(short))
	for i := 0; i < len(long); i++ {
		result[i] = long[i] ^ short[i]
	}
	return string(result)
}
