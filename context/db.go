package context

import (
	"context"
	"github.com/unectio/util/mongo"
)

var ctxSession mongo.Session

func GetDb(c context.Context) mongo.Session {
	return c.(*Context).S
}

func SetDB(c context.Context, s mongo.Session) {
	if ctxSession != nil {
		ctxSession.Close()
	}
	ctxSession = s.Copy()

	ctx := c.(*Context)
	if ctx.S != nil {
		ctx.S.Close()
	}
	ctx.S = s.Copy()
}
