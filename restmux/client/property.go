package client

import rq "github.com/unectio/util/request"

type Property struct {
	Name string
	Col  *Collection
}

func (p *Property) Get(oid string) *rq.Request {
	return rq.Req("", p.Col.Name+"/"+oid+"/"+p.Name).M("GET")
}

func (p *Property) Set(oid string, v interface{}) *rq.Request {
	return rq.Req("", p.Col.Name+"/"+oid+"/"+p.Name).M("PUT").B(v)
}
