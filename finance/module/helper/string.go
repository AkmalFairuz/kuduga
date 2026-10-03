package helper

import "strconv"

func StringRemoveNonAlphanumeric(s string) string {
	var ret string
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') {
			ret += string(c)
		}
	}
	return ret
}

func StringRemoveNonAlphanumericWithoutSpace(s string) string {
	var ret string
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || c == ' ' {
			ret += string(c)
		}
	}
	return ret
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
