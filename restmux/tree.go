package restmux

import (
	"context"
	"net/http"
)

type c struct {
	c	Collection
	e	Error
}

func C(ctx context.Context, x Collection) *c {
	e := x.Acc(ctx)
	if e != nil {
		return &c{ e: e }
	}

	return &c{ c: x }
}

func (c *c)O(ctx context.Context, r *http.Request) *o {
	if c.e != nil {
		return &o{ e: c.e }
	}

	x, err := findObject(ctx, r, c.c)
	if err != nil {
		return &o{ e: err }
	}

	return &o{ c: c.c, o: x }
}

type o struct {
	c	Collection
	o	Object
	e	Error
}

func (o *o)C(x string) *c {
	if o.e != nil {
		return &c{ e: o.e }
	}

	col := o.o.Col(x)
	if col != nil {
		return &c{ c: o.o.Col(x) }
	} else {
		return &c{ e: NotImplementedErr }
	}
}

func (o *o)P(x Property) *p {
	if o.e != nil {
		return &p{ e: o.e }
	}

	return &p{ o: o.o, p: x }
}

func (o *o)A(x Action) *a {
	if o.e != nil {
		return &a{ e: o.e }
	}

	return &a{ o: o.o, a: x }
}

func (o *o)Get() (Object, Error) {
	return o.o, o.e
}

type p struct {
	o	Object
	p	Property
	e	Error
}

type a struct {
	o	Object
	a	Action
	e	Error
}
