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
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestLocal(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("this test can be reproduce only on Linux platform")
	}
	f, err := OpenSafe(".", "vfs_test.go")
	if err != nil {
		t.Fatalf("err: %s\n", err.Error())
	}

	f.Close()

	f, err = OpenSafe(".", "./vfs.go")
	if err == nil {
		fmt.Printf("opened %s\n", f.Name())
		f.Close()
		t.FailNow()
	}

	f, err = OpenSafe(".", "./vfs_test.go")
	if err != nil {
		t.Fatalf("err: %s\n", err.Error())
	}

	f.Close()
}

func TestRmdirAsync(t *testing.T) {
	err := os.Mkdir("./.test_to_remove", 0700)
	if err != nil {
		t.Fatalf("Failed to make dir: %s", err.Error())
	}
	err = Rmdir(".", ".test_to_remove")
	if err != nil {
		t.Fatalf("Failed to remove: %s", err.Error())
	}
}

func chkPath(path string, s bool, t *testing.T) {
	if StraightPath(path) != s {
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

type TestEntry struct {
	DEntry `json:"-"`

	Type  string       `json:"type"`
	Name  string       `json:"name,omitempty"`
	SPath string       `json:"path,omitempty"`
	Kids  []*TestEntry `json:"kids,omitempty"`
}

func makeRoot() *TestEntry {
	return &TestEntry{
		Type:  "dir",
		Name:  "/",
		SPath: "",
	}
}

func buildTree(ent os.FileInfo, dir DEntry) DEntry {
	e := TestEntry{}
	if ent.IsDir() {
		if ent.Name() == ".git" {
			return nil
		}

		e.Type = "dir"
	} else {
		e.Type = "file"
	}
	de := dir.(*TestEntry)
	e.Name = ent.Name()
	e.SPath = de.SPath + "/" + e.Name
	de.Kids = append(de.Kids, &e)
	return &e
}

func buildList(ent os.FileInfo, dir DEntry) (*TestEntry, bool) {
	e := TestEntry{}
	e.SPath = dir.Path() + "/" + ent.Name()

	f := false
	if !ent.IsDir() {
		e.Type = "file"
		f = true
	}

	return &e, f
}

func (e *TestEntry) Path() string { return e.SPath }

func TestWalkTree(t *testing.T) {
	root := makeRoot()
	err := WalkTree("..", root, buildTree)
	if err != nil {
		fmt.Printf("Error build tree: %s\n", err.Error())
		return
	}

	d, err := json.MarshalIndent(root, "", "    ")
	if err != nil {
		fmt.Printf("Error mashalling tree: %s\n", err.Error())
		return
	}

	fmt.Printf("TREE:====================\n%s\n", string(d))
}

func TestWalkList(t *testing.T) {
	root := makeRoot()
	list := []*TestEntry{}
	err := WalkTree("..", root, func(e os.FileInfo, d DEntry) DEntry {
		if e.Name() == ".git" {
			return nil
		}

		ent, f := buildList(e, d)
		if f {
			list = append(list, ent)
		}
		return ent
	})
	if err != nil {
		fmt.Printf("Error build tree: %s\n", err.Error())
		return
	}

	d, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		fmt.Printf("Error mashalling tree: %s\n", err.Error())
		return
	}

	fmt.Printf("LIST:====================\n%s\n", string(d))
}
