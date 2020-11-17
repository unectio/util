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
	"errors"
	"testing"
	"time"
)

func TestCachedValue(t *testing.T) {
	var called bool

	cv := MakeCachedValue(func() (interface{}, error) {
		if called {
			return nil, errors.New("already")
		}

		called = true
		return "foo", nil
	})

	for i := 0; i < 3; i++ {
		x, err := cv.Get()
		if err != nil {
			t.FailNow()
			return
		}

		if x.(string) != "foo" {
			t.FailNow()
			return
		}
	}
}

func TestInvalidate(t *testing.T) {
	var called int
	var cv *CachedValue

	cv = MakeCachedValue(func() (interface{}, error) {
		called++
		time.AfterFunc(500*time.Millisecond, func() { cv.Invalidate() })
		return "foo", nil
	})

	for i := 0; i < 3; i++ {
		cv.Get() //nolint:errcheck
		cv.Get() //nolint:errcheck
		time.Sleep(time.Second)
	}

	if called != 3 {
		t.Fail()
	}
}
