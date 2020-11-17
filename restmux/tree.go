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
	"net/http"
)

type c struct {
	c Collection
	e Error
}

func C(ctx context.Context, x Collection) *c {
	return &c{c: x}
}

func (c *c) O(ctx context.Context, r *http.Request) *o {
	if c.e != nil {
		return &o{e: c.e}
	}

	x, err := findObject(ctx, r, c.c)
	if err != nil {
		return &o{e: err}
	}

	return &o{c: c.c, o: x}
}

type o struct {
	c Collection
	o Object
	e Error
}

func (o *o) C(x string) *c {
	if o.e != nil {
		return &c{e: o.e}
	}

	col := o.o.Col(x)
	if col != nil {
		return &c{c: o.o.Col(x)}
	} else {
		return &c{e: NotImplementedErr}
	}
}

func (o *o) P(x Property) *p {
	if o.e != nil {
		return &p{e: o.e}
	}

	return &p{c: o.c, o: o.o, p: x}
}

func (o *o) A(x Action) *a {
	if o.e != nil {
		return &a{e: o.e}
	}

	return &a{c: o.c, o: o.o, a: x}
}

func (o *o) Get() (Object, Error) {
	return o.o, o.e
}

type p struct {
	c Collection
	o Object
	p Property
	e Error
}

type a struct {
	c Collection
	o Object
	a Action
	e Error
}
