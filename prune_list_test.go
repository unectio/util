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
	"testing"
	"time"
)

func TestPruneList(t *testing.T) {
	pl := NewPruneList()

	/* Double sched */
	pl.Schedule("1", "one", 3*time.Second)
	pl.Schedule("1", "one", 3*time.Second)

	/* Unsched */
	pl.Schedule("2", "two", 10*time.Second)
	pl.Schedule("2", "two", 10*time.Second)
	pl.Unschedule("2")

	/* Resched */
	pl.Schedule("3", "three", 5*time.Second)
	pl.Unschedule("3")
	pl.Schedule("3", "three", 5*time.Second)

	/* Schedule earlier */
	pl.Schedule("4", "four", 2*time.Second)

	x := <-pl.Wait()
	fmt.Printf("%s\n", x.(string))

	x = <-pl.Wait()
	fmt.Printf("%s\n", x.(string))

	x = <-pl.Wait()
	fmt.Printf("%s\n", x.(string))
}
