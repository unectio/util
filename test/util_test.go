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
	"os"
	"fmt"
	"time"
	"testing"
	"github.com/unectio/util"
	"github.com/unectio/util/sync"
)

func TestCreds(t *testing.T) {
	usr := "foo"
	pwd := "1q2w3e"
	adr := "host.123.com"
	prt := "12345"
	dom := "domain"

	c := util.CredsParse(usr + ":" + pwd + "@" + adr + ":" + prt + "/" + dom)
	if c.Usr != usr || c.Pwd != pwd || c.Adr != adr || c.Prt != prt || c.Dom != dom {
		t.Fatalf("Unexpected %v", c)
	}
}

func TestStrinc(t *testing.T) {
	v := util.StrintInc("42")
	if v != "43" {
		t.Fatalf("%s != 43", v)
	}
}

func TestRandstr(t *testing.T) {
	s, _ := util.RandomString(12)
	if len(s) != 12 {
		t.Fatalf("len(%s) == %d", s, len(s))
	}
	for i, r := range(s) {
		if !(r == '_' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			t.Fatalf("bad chr @%d in %s", i, s)
		}
	}
}

func TestRmdirAsync(t *testing.T) {
	err := os.Mkdir("./.test_to_remove", 0700)
	if err != nil {
		t.Fatalf("Failed to make dir: %s", err.Error())
	}
	err = util.Rmdir(".", ".test_to_remove")
	if err != nil {
		t.Fatalf("Failed to remove: %s", err.Error())
	}
}

func chkPath(path string, s bool, t *testing.T) {
	if util.StraightPath(path) != s {
		t.Errorf("%s != %v\n", path, s)
		t.Fail()
	}
}

func TestStraighPath(t *testing.T) {
	chkPath("..", false, t)
	chkPath("/..", false, t)
	chkPath("../", false, t)
	chkPath("foo/..", false, t)
	chkPath("../foo", false, t)
	chkPath("foo/../foo", false, t)
	chkPath("../foo/..", false, t)

	chkPath(".", true, t)
	chkPath("/.", true, t)
	chkPath("./", true, t)
	chkPath("foo/.", true, t)
	chkPath("./foo", true, t)
	chkPath("foo/./foo", true, t)
	chkPath("./foo/.", true, t)

	chkPath("...", true, t)
	chkPath("/...", true, t)
	chkPath(".../", true, t)
	chkPath("foo/...", true, t)
	chkPath(".../foo", true, t)
	chkPath("foo/.../foo", true, t)
	chkPath(".../foo/...", true, t)

	chkPath("foo.bar", true, t)
	chkPath("foo..bar", true, t)
	chkPath("foo...bar", true, t)

	chkPath(".foo", true, t)
	chkPath("..foo", true, t)
	chkPath("...foo", true, t)
}

func TestPint(t *testing.T) {
	pint := util.PInt{}
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

func TestTill(t *testing.T) {
	fmt.Printf("%s Start\n", time.Now())
	sync.TillMinute.Sleep()
	fmt.Printf("%s Stop\n", time.Now())
}

func TestMergeUniq(t *testing.T) {
	fmt.Printf("%v\n", util.MergeUniq([]string{"foo", "bar"}, []string{"foo", "buzz"}))
}
