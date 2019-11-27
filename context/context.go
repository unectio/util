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
