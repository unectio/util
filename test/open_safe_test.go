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
	"github.com/unectio/util"
)

func TestLocal(t *testing.T) {
	f, err := util.OpenSafe(".", "open_safe_test.go")
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		t.FailNow()
	}

	f.Close()

	f, err = util.OpenSafe(".", "../vfs.go")
	if err == nil {
		fmt.Printf("opened %s\n", f.Name())
		f.Close()
		t.FailNow()
	}

	f, err = util.OpenSafe(".", "../test/open_safe_test.go")
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		t.FailNow()
	}

	f.Close()
}

