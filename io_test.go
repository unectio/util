//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package util

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestReadCmdLines(t *testing.T) {
	cmd := exec.Command("ls", ".")
	l, err := ReadCmdLines(cmd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", strings.Join(l, "\n"))
}
