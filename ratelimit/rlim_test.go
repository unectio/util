//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

var (
	burst uint = 10
	rate  uint = 5
)

func TestRatelimit(t *testing.T) {
	rl := NewFilter(burst, rate)

	stop := time.Now().Add(time.Second)
	events := uint(0)
	for {
		if rl.Step() {
			events++
		}
		time.Sleep(time.Millisecond)
		if time.Now().After(stop) {
			break
		}
	}

	exp := burst + rate
	fmt.Printf("%d/%d events\n", events, exp)
	if events > 2*exp {
		t.Fail()
	}
}
