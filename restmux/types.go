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
	Find(context.Context /* id */, string) (Object, Error)
	Lookup(context.Context /* name */, string, url.Values) (Object, Error)
	Add(context.Context, Image) (Object, Error)
	Del(context.Context, Object) Error
	Iter(context.Context, url.Values, func(context.Context, Object) Error) Error

	Image() Image
	/*
	 * Parameter name that's used in URL to 'encode' an object ID
	 */
	Id() string
}
