package context

import (
	"net/http"
)

type LoginInfo interface {
	Scope() string
}

type simpleLogin string

func (sl simpleLogin)Scope() string { return string(sl) }

type HttpLogin struct {
	simpleLogin
}

func GetHttpLogin(r *http.Request) HttpLogin {
	return HttpLogin{
		simpleLogin(r.RemoteAddr),
	}
}

func GetHttpLogin2(r *http.Request, pfx string) HttpLogin {
	return HttpLogin{
		simpleLogin(pfx + ":" + r.RemoteAddr),
	}
}

type TestLogin struct {
	simpleLogin
}

func GetTestLogin() TestLogin {
	return TestLogin{
		simpleLogin("test"),
	}
}
