/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

package context

import (
	"context"
	"sync/atomic"
	"github.com/unectio/util/mongo"
	"go.uber.org/zap"
)

var rover uint64

type Context struct {
	context.Context
	S	mongo.Session
	L	*zap.SugaredLogger

	rover	uint64
	li	LoginInfo
}

func Make(li LoginInfo) *Context {
	ctx := Context{}
	ctx.rover = atomic.AddUint64(&rover, 1)
	ctx.li = li
	ctx.init()
	return &ctx
}

type forkedLogin struct {
	from	LoginInfo
	reason	string
}

func (fl *forkedLogin)Scope() string {
	return fl.from.Scope() + " (" + fl.reason + ")"
}

func Fork(ctx context.Context, reason string) *Context {
	pctx := ctx.(*Context)
	nctx := Context{}
	nctx.rover = pctx.rover
	nctx.li = &forkedLogin{pctx.li, reason}
	nctx.init()
	return &nctx
}

func Login(ctx context.Context) LoginInfo {
	return ctx.(*Context).li
}

func L(ctx context.Context) *zap.SugaredLogger {
	return ctx.(*Context).L
}

func (ctx *Context) init() {
	ctx.Context = context.Background()
	ctx.L = dlog.With(zap.Int64("r", int64(ctx.rover)), zap.String("s", ctx.li.Scope()))
	if ctxSession != nil {
		ctx.S = ctxSession.Copy()
	}
}

func (ctx *Context) Close() {
	if ctx.S != nil {
		ctx.S.Close()
		ctx.S = nil
	}
}
