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

package util

import (
	"io"
	"time"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	BearerPrefix string = "Bearer "
)

func ReadJsonBody(body io.ReadCloser, into interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(into)
}

func QueryParam(q url.Values, name string) string {
	vals := q[name]
	if len(vals) > 0 {
		return vals[0]
	} else {
		return ""
	}
}

func QueryDuration(q url.Values, name string) (bool, time.Duration) {
	p := QueryParam(q, name)
	if p == "" {
		return true, 0
	}

	d, err := time.ParseDuration(p)
	if err == nil {
		return true, d
	}

	return false, 0
}

func UrlQuery(name, value string) url.Values {
	return url.Values{name: []string{value}}
}

func ParseBearer(token string) (string, bool) {
	return CutPrefix(token, BearerPrefix)
}

func QueryArgs(r *http.Request) map[string]string {
	ret := make(map[string]string)

	for name, value := range r.URL.Query() {
		if len(value) == 1 {
			ret[name] = value[0]
		}
	}

	return ret
}

func QueryBody(r *http.Request) (string, []byte) {
	ct := r.Header.Get("Content-Type")
	aux := strings.SplitN(ct, ";", 2)
	if len(aux) > 0 {
		switch (aux[0]) {
			case "text/plain":
			case "application/json":
				body, err := ioutil.ReadAll(r.Body)
				if err == nil && len(body) > 0 {
					return aux[1], body
				}
		}
	}

	return "", nil
}
