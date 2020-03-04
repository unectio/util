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
	"testing"
	"github.com/unectio/util/seqpack"
)

type Foo struct {
	Bar	int
}

func TestSeqPack(t *testing.T) {
	sk, err := seqpack.Make()
	if err != nil {
		t.Fatalf("Cannot make socket: %s", err.Error())
	}

	go func() {
		psk, err := seqpack.Open(sk.PFd())
		if err != nil {
			t.Fatalf("Cannot open peer sk: %s", err.Error())
		}

		r1 := map[string]string{}
		err = psk.Recv(&r1)
		if err != nil {
			t.Fatalf("Cannot recv map: %s", err.Error())
		}

		fmt.Printf("1: %v\n", r1)

		err = psk.Send(&Foo{42})
		if err != nil {
			t.Fatalf("Cannot recv struct: %s", err.Error())
		}
	}()

	err = sk.Send(map[string]string{"foo":"bar"})
	if err != nil {
		t.Fatalf("Cannot send map: %s", err.Error())
	}

	r := &Foo{}
	err = sk.Recv(r)
	if err != nil {
		t.Fatalf("Cannot recv struct: %s", err.Error())
	}

	fmt.Printf("2: %v\n", r)
}
