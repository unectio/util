package router

import (
	"errors"
	"strings"
)

type Router struct {
	root *layer
}

func MakeRouter() *Router {
	rt := &Router{}
	rt.root = new_layer()
	return rt
}

func (r *Router) RegisterURL(url string, res interface{}) error {
	debug("+ URL [%s]\n", url)
	path := strings.Split(url, "/")
	return r.root.register_path(path, res)
}

func (r *Router) HandleURL(url string) (interface{}, map[string]string) {
	debug("= URL [%s]\n", url)
	path := strings.Split(url, "/")
	return r.root.resolve_path(path)
}

type layer struct {
	exact map[string]*layer
	param map[string]*layer
	match interface{}
}

func new_layer() *layer {
	pm := &layer{}
	pm.exact = make(map[string]*layer)
	pm.param = make(map[string]*layer)
	return pm
}

func is_param(n string) (bool, string) {
	if strings.HasPrefix(n, "{") && strings.HasSuffix(n, "}") {
		return true, n[1 : len(n)-1]
	} else {
		return false, ""
	}
}

func next_layer(lrs map[string]*layer, name string) *layer {
	next, ok := lrs[name]
	if !ok {
		next = new_layer()
		lrs[name] = next
	}

	return next
}

func (lr *layer) register_path(path []string, res interface{}) error {
	debug("\t+ PATH [%v]\n", path)
	if len(path) == 0 {
		if lr.match != nil {
			return errors.New("name conflict")
		}

		lr.match = res
		return nil
	}

	cur := path[0]
	var next *layer

	if ok, param := is_param(cur); ok {
		next = next_layer(lr.param, param)
	} else {
		next = next_layer(lr.exact, cur)
	}

	return next.register_path(path[1:], res)
}

func (lr *layer) resolve_path(path []string) (interface{}, map[string]string) {
	debug("\t= PATH [%v]\n", path)
	if len(path) == 0 {
		return lr.match, nil
	}

	cur := path[0]

	next, ok := lr.exact[cur]
	if ok {
		return next.resolve_path(path[1:])
	}

	for pn, p := range lr.param {
		res, params := p.resolve_path(path[1:])
		if res != nil {
			if params == nil {
				params = make(map[string]string)
			}
			params[pn] = cur
			return res, params
		}
	}

	return nil, nil
}
