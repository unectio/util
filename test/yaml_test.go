package main

import (
	"testing"
	"github.com/unectio/util"
)

func check(foo *util.YAMLRaw) {
}

func TestRouterUnmarshal(t *testing.T) {
	var foo util.YAMLRaw
	check(&foo)
}

