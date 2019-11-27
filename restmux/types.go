package restmux

import (
	"context"
	"net/url"
)

type Image interface{}

type Object interface {
	Init(context.Context, Collection)
	Info(context.Context, url.Values, bool) Image
	Upd(context.Context, url.Values, Image) Error

	Col(string) Collection
}

type Property interface {
	Get(context.Context, Object, url.Values) (Image, Error)
	Set(context.Context, Object, Image) Error
	Del(context.Context, Object) Error

	Image() Image
}

type Action interface {
	Do(context.Context, Object, Image, url.Values) (Image, Error)

	Image() Image
}

type Collection interface {
	Acc(context.Context) Error
	Find(context.Context, /* id */ string) (Object, Error)
	Lookup(context.Context, /* name */ string, url.Values) (Object, Error)
	Add(context.Context, Image) (Object, Error)
	Del(context.Context, Object) Error
	Iter(context.Context, url.Values, func(context.Context, Object) Error) Error

	Image() Image
	/*
	 * Parameter name that's used in URL to 'encode' an object ID
	 */
	Id() string
}
