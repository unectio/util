//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"testing"
	"strings"
	"os/exec"
	"github.com/unectio/util"
)

func TestReadCmdLines(t *testing.T) {
	cmd := exec.Command("ls", ".")
	l, err := util.ReadCmdLines(cmd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", strings.Join(l, "\n"))
}
