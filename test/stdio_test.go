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
	"time"
	"testing"
	"github.com/unectio/util/stdio"
)

func TestStream(t *testing.T) {
	out, err := stdio.Make("out")
	if err != nil {
		fmt.Printf("error making stream: %s\n", err.Error())
		t.FailNow()
		return
	}

	tmo := time.AfterFunc(5 * time.Second, func() {
					fmt.Printf("Timeout\n")
					t.FailNow()
				})

	x := out.Read()
	tmo.Stop()
	fmt.Printf("[%s]\n", x)
}
