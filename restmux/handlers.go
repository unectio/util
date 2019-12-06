package restmux

import (
	"context"
	"net/url"
	"net/http"
	"github.com/gorilla/mux"
)

var errBadm = &GenError{http.StatusMethodNotAllowed, ""}

func resolve(ctx context.Context, w http.ResponseWriter, name string, q url.Values, c Collection) Error {
	o, err := c.Lookup(ctx, name, q)
	if err != nil {
		return err
	}

	o.Init(ctx, c)

	return respondJson(ctx, w, o.Info(ctx, nil, false))
}

func list(ctx context.Context, w http.ResponseWriter, r *http.Request, c Collection) Error {
	q := r.URL.Query()
	name := q.Get("name")
	if name != "" {
		return resolve(ctx, w, name, q, c)
	}

	ret := make([]Image, 0) /* Ensure [], not nil in marshal */
	details := (q.Get("details") != "")

	err := c.Iter(ctx, q, func(ctx context.Context, o Object) Error {
		ret = append(ret, o.Info(ctx, q, details))
		return nil
	})

	if err != nil {
		return err
	}

	return respondJson(ctx, w, ret)
}

func create(ctx context.Context, w http.ResponseWriter, r *http.Request, c Collection) Error {
	i := c.Image()
	err := read(r.Body, i)
	if err != nil {
		return err
	}

	o, err := c.Add(ctx, i)
	if err != nil {
		return err
	}

	return respondJson2(ctx, w, http.StatusCreated, o.Info(ctx, nil, false))
}

func (c *c)Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) Error {
	if c.e != nil {
		return c.e
	}

	e := c.c.Acc(ctx)
	if e != nil {
		return e;
	}

	switch r.Method {
	case "GET":
		return list(ctx, w, r, c.c)
	case "POST":
		return create(ctx, w, r, c.c)
	}

	return errBadm
}

var CMethods = []string{ "GET", "POST", "OPTIONS" }

func info(ctx context.Context, w http.ResponseWriter, r *http.Request, o Object) Error {
	return respondJson(ctx, w, o.Info(ctx, r.URL.Query(), true))
}

func update(ctx context.Context, w http.ResponseWriter, r *http.Request, c Collection, o Object) Error {
	i := c.Image()
	err := read(r.Body, i)
	if err != nil {
		return err
	}

	err = o.Upd(ctx, r.URL.Query(), i)
	if err != nil {
		return err
	}

	return respondJson(ctx, w, o.Info(ctx, nil, true))
}

func remove(ctx context.Context, w http.ResponseWriter, r *http.Request, c Collection, o Object) Error {
	err := c.Del(ctx, o)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func findObject(ctx context.Context, r *http.Request, c Collection) (Object, Error) {
	o, err := c.Find(ctx, mux.Vars(r)[c.Id()])
	if err == nil {
		o.Init(ctx, c)
	}
	return o, err
}

func (o *o)Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) Error {
	if o.e != nil {
		return o.e
	}

	e := o.c.Acc(ctx)
	if e != nil {
		return e;
	}

	switch r.Method {
	case "GET":
		return info(ctx, w, r, o.o)
	case "PUT":
		return update(ctx, w, r, o.c, o.o)
	case "DELETE":
		return remove(ctx, w, r, o.c, o.o)
	}

	return errBadm
}

var OMethods = []string{ "GET", "PUT", "DELETE", "OPTIONS" }

func pget(ctx context.Context, w http.ResponseWriter, r *http.Request, o Object, prop Property) Error {
	i, err := prop.Get(ctx, o, r.URL.Query())
	if err != nil {
		return err
	}

	return respondJson(ctx, w, i)
}

func pset(ctx context.Context, w http.ResponseWriter, r *http.Request, o Object, prop Property) Error {
	i := prop.Image()
	if i != nil {
		err := read(r.Body, i)
		if err != nil {
			return err
		}
	}

	err := prop.Set(ctx, o, i)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func pdel(ctx context.Context, w http.ResponseWriter, r *http.Request, o Object, prop Property) Error {
	err := prop.Del(ctx, o)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (p *p)Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) Error {
	if p.e != nil {
		return p.e
	}

	e := p.c.Acc(ctx)
	if e != nil {
		return e;
	}

	switch r.Method {
	case "GET":
		return pget(ctx, w, r, p.o, p.p)
	case "PUT":
		return pset(ctx, w, r, p.o, p.p)
	case "DELETE":
		return pdel(ctx, w, r, p.o, p.p)
	}

	return errBadm
}

var PMethods = []string{ "GET", "PUT", "DELETE", "OPTIONS" }

func act(ctx context.Context, w http.ResponseWriter, r *http.Request, o Object, a Action) Error {
	i := a.Image()
	if i != nil {
		err := read(r.Body, i)
		if err != nil {
			return err
		}
	}

	i, err := a.Do(ctx, o, i, r.URL.Query())
	if err != nil {
		return err
	}

	if i != nil {
		respondJson(ctx, w, i)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return nil
}

func (a *a)Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) Error {
	if a.e != nil {
		return a.e
	}

	e := a.c.Acc(ctx)
	if e != nil {
		return e;
	}

	switch r.Method {
	case "POST":
		return act(ctx, w, r, a.o, a.a)
	}

	return errBadm
}

var AMethods = []string{ "POST", "OPTIONS" }
