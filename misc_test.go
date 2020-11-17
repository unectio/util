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
	"testing"
)

func TestStrinc(t *testing.T) {
	v := StrintInc("42")
	if v != "43" {
		t.Fatalf("%s != 43", v)
	}
}

func TestMergeUniq(t *testing.T) {
	fmt.Printf("%v\n", MergeUniq([]string{"foo", "bar"}, []string{"foo", "buzz"}))
}
