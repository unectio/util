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

func (r *Response)Hdrs() string {
	ret := ""
	for h, v := range r.resp.Header {
		ret += fmt.Sprintf("%s=%s;", h, v)
	}

	return ret
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
