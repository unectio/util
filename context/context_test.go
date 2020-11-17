//////////////////////////////////////////////////////////////////////////////
//
// (C) Copyright 2019-2020 by Unectio, Inc.
//
// The information contained herein is confidential, proprietary to Unectio,
// Inc.
//
//////////////////////////////////////////////////////////////////////////////

package context

import (
	"testing"
)

func TestContextLog(t *testing.T) {
	ctx := Make(Global("test"))
	L(ctx).Infof("Hello, world!")
}

func TestContextFork(t *testing.T) {
	ctx := Make(Global("test-fork"))
	ctx2 := Fork(ctx, "test")
	L(ctx).Infof("Hello, kid!")
	L(ctx2).Infof("Hello, fork!")
}
