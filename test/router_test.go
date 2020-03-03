package main

import (
	"fmt"
	"testing"
	"github.com/unectio/util/restmux/router"
)

func TestRouter(t *testing.T) {
	r := router.MakeRouter()

	router.Debug()

	if r.RegisterURL("foo", 1) != nil {
		fmt.Printf("Name conflict foo")
		t.FailNow()
	}
	r.Print()

	if r.RegisterURL("foo/{id}", 2) != nil {
		fmt.Printf("Name conflict foo/{id}")
		t.FailNow()
	}
	r.Print()

	if r.RegisterURL("foo/bar", 3) != nil {
		fmt.Printf("Name conflict foo/bar")
		t.FailNow()
	}
	r.Print()

	if r.RegisterURL("foo/{id}/bar", 4) != nil {
		fmt.Printf("Name conflict foo/{id}/bar")
		t.FailNow()
	}
	r.Print()

	if r.RegisterURL("foo/{id}/bar/{id2}", 5) != nil {
		fmt.Printf("Name conflict foo/{id}/bar/{id2}")
		t.FailNow()
	}
	r.Print()

	if r.RegisterURL("foo/{id}/{id2}", 6) != nil {
		fmt.Printf("Name conflict foo/{id}/{id2}")
		t.FailNow()
	}
	r.Print()


	res, p := r.HandleURL("foo")
	if res != 1 || p != nil {
		fmt.Printf("foo -> %v/%v\n", res, p)
		t.Fail()
	}

	res, p = r.HandleURL("foo/FOO")
	if res != 2 || p["id"] != "FOO" || len(p) != 1 {
		fmt.Printf("foo/FOO -> %v/%v\n", res, p)
		t.Fail()
	}

	res, p = r.HandleURL("foo/bar")
	if res != 3 || p != nil {
		fmt.Printf("foo/bar -> %v/%v\n", res, p)
		t.Fail()
	}

	res, p = r.HandleURL("foo/FOO/bar")
	if res != 4 || p["id"] != "FOO" || len(p) != 1 {
		fmt.Printf("foo/FOO/bar -> %v/%v\n", res, p)
		t.Fail()
	}

	res, p = r.HandleURL("foo/FOO/bar/BAR")
	if res != 5 || p["id"] != "FOO" || p["id2"] != "BAR" || len(p) != 2 {
		fmt.Printf("foo/FOO/bar/BAR -> %v/%v\n", res, p)
		t.Fail()
	}

	res, p = r.HandleURL("foo/FOO/BAR")
	if res != 6 || p["id"] != "FOO" || p["id2"] != "BAR" || len(p) != 2 {
		fmt.Printf("foo/FOO/BAR -> %v/%v\n", res, p)
		t.Fail()
	}
}
