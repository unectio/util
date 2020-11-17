//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package sync

import (
	"fmt"
	"testing"
	"time"
)

func TestTill(t *testing.T) {
	fmt.Printf("%s Start\n", time.Now())
	// TODO: check that it works without sleep
	// TillMinute.Sleep()
	fmt.Printf("%s Stop\n", time.Now())
}
