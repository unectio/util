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

func TestCreds(t *testing.T) {
	usr := "foo"
	pwd := "1q2w3e"
	adr := "host.123.com"
	prt := "12345"
	dom := "domain"

	c := CredsParse(usr + ":" + pwd + "@" + adr + ":" + prt + "/" + dom)
	if c.Usr != usr || c.Pwd != pwd || c.Adr != adr || c.Prt != prt || c.Dom != dom {
		t.Fatalf("Unexpected %v", c)
	}
}
