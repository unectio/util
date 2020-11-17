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

func TestPint(t *testing.T) {
	pint := PInt{}
	pint.Inc()
	pint.Inc()
	pint.Inc()
	if pint.Peak != 3 {
		t.Fail()
	}
	pint.Dec()
	if pint.Peak != 3 {
		t.Fail()
	}
	pint.Dec()
	pint.Inc()
	if pint.Peak != 3 {
		t.Fail()
	}
}
