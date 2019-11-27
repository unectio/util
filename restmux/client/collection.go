package client

import (
	"net/http"
	"github.com/unectio/util/request"
)

type Collection struct {
	Name	string
	Parent	*Collection
}

func (c *Collection)List() *rq.Request {
	return rq.Req("", c.Name).M("GET")
}

func (c *Collection)Add(o interface{}) *rq.Request {
	return rq.Req("", c.Name).B(o).OK(http.StatusCreated)
}

func (c *Collection)Info(id string) *rq.Request {
	return rq.Req("", c.Name + "/" + id).M("GET")
}

func (c *Collection)Upd(id string, o interface{}) *rq.Request {
	return rq.Req("", c.Name + "/" + id).M("PUT").B(o)
}

func (c *Collection)Lookup(n string) *rq.Request {
	return rq.Req("", c.Name + "?name=" + n).M("GET")
}

func (c *Collection)Delete(id string) *rq.Request {
	return rq.Req("", c.Name + "/" + id).M("DELETE").OK(http.StatusNoContent)
}

func (c *Collection)Sub(pid string) *Collection {
	return &Collection {
		Name:	c.Parent.Name + "/" + pid + "/" + c.Name,
	}
}
