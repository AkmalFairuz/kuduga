package helper

import (
	"strings"
)

func IPVersion(ip string) int {
	if strings.Contains(ip, ".") {
		return 4
	}
	if strings.Contains(ip, ":") {
		return 6
	}
	return 0
}

func GetIPV6Prefix(ip string) string {
	if IPVersion(ip) != 6 {
		return ""
	}
	parts := strings.Split(ip, ":")
	l := min(4, len(parts))
	ret := make([]string, 0, l)
	for i := 0; i < l; i++ {
		ret = append(ret, parts[i])
	}
	return strings.Join(ret, ":")
}
