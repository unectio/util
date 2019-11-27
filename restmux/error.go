package restmux

import (
	"fmt"
	"net/http"
)

type Error interface {
	String() string
	Status() int
}

type GenError struct {
	Code		int
	Message		string
}

func (e *GenError)String() string {
	if e.Message == "" {
		return fmt.Sprintf("status=%d", e.Code)
	}

	return e.Message
}

func (e *GenError)Status() int {
	return e.Code
}

func HttpErr(w http.ResponseWriter, e Error) {
	http.Error(w, e.String(), e.Status())
}

var NotImplementedErr = &GenError{http.StatusMethodNotAllowed, "not implemented"}

type BadReqErr string

func (e BadReqErr)String() string {
	return string(e)
}

func (e BadReqErr)Status() int {
	return http.StatusBadRequest
}

func BadReq(err error) BadReqErr {
	return BadReqErr(err.Error())
}
