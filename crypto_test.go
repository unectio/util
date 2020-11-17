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
	"testing"
)

func TestRandstr(t *testing.T) {
	s, _ := RandomString(12)
	if len(s) != 12 {
		t.Fatalf("len(%s) == %d", s, len(s))
	}
	for i, r := range s {
		if !(r == '_' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			t.Fatalf("bad chr @%d in %s", i, s)
		}
	}
}
