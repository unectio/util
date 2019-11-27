package context

import (
	"net/http"
)

type LoginInfo interface {
	Scope() string
}

type HttpLogin string

func (sc HttpLogin)Scope() string {
	return string(sc)
}

func GetHttpLogin(r *http.Request) HttpLogin {
	return HttpLogin(r.RemoteAddr)
}

func GetHttpLogin2(r *http.Request, pfx string) HttpLogin {
	return HttpLogin(pfx + ":" + r.RemoteAddr)
}
