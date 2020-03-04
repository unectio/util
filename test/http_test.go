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
	"time"
	"testing"
	"net/http"
	"github.com/unectio/util/request"
)

func server(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("Path:   %s\n", r.URL.Path)
	h := r.Header.Get("X-Test-Head")
	fmt.Printf("Header: %s\n", h)
	q := append(r.URL.Query()["query"], "")
	fmt.Printf("Query:  %s\n", q)

	if r.Method == "PATCH" && r.URL.Path == "/foo/bar" && h == "yes" && q[0] == "param" {
		w.WriteHeader(http.StatusGone)
	} else {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func TestRestReq(t *testing.T) {
	server := http.Server{
		Handler:      http.HandlerFunc(server),
		Addr:         "127.0.0.1:8888",
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			t.Fatalf("err in server: %s", err.Error())
		}
	}()

	/* How else to wait for the server port open? */
	time.Sleep(500 * time.Millisecond)

	resp := rq.Req("http://127.0.0.1:8888", "/foo/bar?query=param").M("PATCH").OK(http.StatusGone).H("X-Test-Head", "yes").Do()
	if !resp.OK() {
		t.Fatalf("err in req: %s", resp.Error())
	}
}
