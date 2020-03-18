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
	"testing"
	"github.com/unectio/util"
	"encoding/json"
)

type TestEntry struct {
	util.DEntry			`json:"-"`

	Type		string		`json:"type"`
	Name		string		`json:"name,omitempty"`
	SPath		string		`json:"path,omitempty"`
	Kids		[]*TestEntry	`json:"kids,omitempty"`
}

func makeRoot() *TestEntry {
	return &TestEntry{
		Type: "dir",
		Name: "/",
		SPath: "",
	}
}

func buildTree(ent os.FileInfo, dir util.DEntry) util.DEntry {
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

func buildList(ent os.FileInfo, dir util.DEntry) (*TestEntry, bool) {
	e := TestEntry{}
	e.SPath = dir.Path() + "/" + ent.Name()

	f := false
	if !ent.IsDir() {
		e.Type = "file"
		f = true
	}

	return &e, f
}

func (e *TestEntry)Path() string { return e.SPath }

func TestWalkTree(t *testing.T) {
	root := makeRoot()
	err := util.WalkTree("..", root, buildTree)
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
	err := util.WalkTree("..", root, func(e os.FileInfo, d util.DEntry) util.DEntry {
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
