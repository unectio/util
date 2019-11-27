package rq

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/unectio/util"
)

type Response struct {
	resp	*http.Response
	err	error
}

func (r *Response)Code() string {
	if r.err == nil {
		return "OK"
	} else if r.resp != nil {
		return "http." + r.resp.Status
	} else {
		return "error"
	}
}

func (r *Response)OK() bool {
	return r.err == nil
}

func (r *Response)E() error {
	if r.OK() {
		return nil
	} else {
		return r
	}
}

func (r *Response)String() string {
	if r.OK() {
		return fmt.Sprintf("OK %s", r.resp.Status)
	} else {
		if r.resp != nil {
			return fmt.Sprintf("%s (status %s)", r.err.Error(), r.resp.Status)
		} else {
			return fmt.Sprintf("%s", r.err.Error())
		}
	}
}

func (r *Response)Error() string {
	if r.OK() {
		return ""
	}

	code := 0
	body := "-"
	if r.resp != nil {
		code = r.resp.StatusCode
		data, _ := ioutil.ReadAll(r.resp.Body)
		body = string(data)
	}

	return fmt.Sprintf("Request failed: %s (response: %d/%s)", r.err.Error(), code, body)
}

func (r *Response)B(out interface{}) *Response {
	if r.OK() {
		err := util.ReadJsonBody(r.resp.Body, out)
		if err != nil {
			r.err = fmt.Errorf("Error decoding response body: %s", err.Error())
		}
	}
	return r
}

func (r *Response)Raw() (*Response, []byte) {
	var body []byte

	if r.OK() {
		var err error

		body, err = ioutil.ReadAll(r.resp.Body)
		if err != nil {
			r.err = fmt.Errorf("Error reading response body: %s", err.Error())
		}
	}

	return r, body
}
