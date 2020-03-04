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
	"testing"
	"github.com/unectio/util/context"
)

func TestContextLog(t *testing.T) {
	ctx := context.Make(context.Global("test"))
	context.L(ctx).Infof("Hello, world!")
}

func TestContextFork(t *testing.T) {
	ctx := context.Make(context.Global("test-fork"))
	ctx2 := context.Fork(ctx, "test")
	context.L(ctx).Infof("Hello, kid!")
	context.L(ctx2).Infof("Hello, fork!")
}
