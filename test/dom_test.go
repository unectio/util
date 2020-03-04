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
	"testing"
	"github.com/unectio/util"
)

func TestDomain(t *testing.T) {
	if !util.IsDomainName("a.com") {
		t.Errorf("a.com not valid")
	}
	if !util.IsDomainName("0.com") {
		t.Errorf("0.com not valid")
	}
	if !util.IsDomainName("a.b.cd") {
		t.Errorf("a.b.cd not valid")
	}
	if !util.IsDomainName("a-b.cd") {
		t.Errorf("a-b.cd not valid")
	}

	if util.IsDomainName("com") {
		t.Errorf("com valid")
	}
	if util.IsDomainName("-a.com") {
		t.Errorf("-a.com valid")
	}
	if util.IsDomainName("a-.com") {
		t.Errorf("a-.com valid")
	}
	if util.IsDomainName("a_z.com") {
		t.Errorf("a_z.com valid")
	}
	if util.IsDomainName(".a.com") {
		t.Errorf(".a.com valid")
	}
	if util.IsDomainName(".com") {
		t.Errorf(".com valid")
	}
	if util.IsDomainName("a.b") {
		t.Errorf("a.b valid")
	}
	if util.IsDomainName("a.b1") {
		t.Errorf("a.b1 valid")
	}

	x := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" /* 64 */
	if util.IsDomainName(x + ".com") {
		t.Errorf(x + ".com valid")
	}
}
