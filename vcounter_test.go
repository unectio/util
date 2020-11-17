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
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestVcount(t *testing.T) {
	vc := MakeVCount()

	vg := sync.WaitGroup{}
	vg.Add(3)
	go func() {
		tmo, _ := vc.Wait("a", time.Second, "=0")
		if tmo {
			t.Errorf("t1")
			t.Fail()
		}
		vg.Done()
	}()
	go func() {
		tmo, _ := vc.Wait("a", time.Second, "=1")
		if !tmo {
			t.Errorf("t2")
			t.Fail()
		}
		vg.Done()
	}()
	go func() {
		tmo, _ := vc.Wait("a", time.Second, "+1")
		if tmo {
			t.Errorf("t3")
			t.Fail()
		}
		vg.Done()
	}()

	time.Sleep(100 * time.Millisecond)
	vc.Add("a", "0")
	vc.Add("a", "2")
	vg.Wait()

	vg = sync.WaitGroup{}
	vg.Add(1)
	go func() {
		tmo, gone := vc.Wait("b", time.Second, "=1")
		if tmo || !gone {
			t.Errorf("t4")
			t.Fail()
		}
		vg.Done()
	}()

	fmt.Printf("Test 'gone' thing\n")
	vc.Add("b", VcDying)
	vg.Wait()
}
