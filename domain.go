package util

import (
	"regexp"
)

var domRe = regexp.MustCompile("^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,63}$")

func IsDomainName(name string) bool {
	if len(name) == 0 || len(name) > 253 {
		return false
	}

	return domRe.MatchString(name)
}
