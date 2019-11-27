package util

import (
	"errors"
	"strings"
	"strconv"
)

func StrintInc(str string) string {
	v, err := strconv.Atoi(str)
	if err != nil {
		return str
	} else {
		return strconv.Itoa(v + 1)
	}
}

func CutPrefix(s, p string) (string, bool) {
	if strings.HasPrefix(s, p) {
		return s[len(p):], true
	} else {
		return s, false
	}
}

func MergeUniq(into, what []string) []string {
	x := make(map[string]struct{})
	for _, k := range into {
		x[k] = struct{}{}
	}
	for _, k := range what {
		if _, ok := x[k]; !ok {
			into = append(into, k)
		}
	}

	return into
}

func Error(msg string, err error) error {
	return errors.New(msg + ": " + err.Error())
}
