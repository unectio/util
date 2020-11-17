package util

import (
	"testing"
)

func check(foo *YAMLRaw) {
}

func TestRouterUnmarshal(t *testing.T) {
	var foo YAMLRaw
	check(&foo)
}
