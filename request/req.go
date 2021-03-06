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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/unectio/util"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	sys_cert_pool string = "::SYS::"
)

type Request struct {
	Method      string
	Host        string
	Path        string
	Headers     map[string]string
	Status      int
	Body        interface{}
	Timeout     time.Duration
	Certificate string

	sign  []byte
	signH string
}

func (rq *Request) URL() string {
	return rq.Host + rq.Path
}

func Req(host, url string) *Request {
	return &Request{
		Host: host,
		Path: url,

		Method: "POST",
		Status: http.StatusOK,
	}
}

func (r *Request) Q(q string) *Request {
	r.Path += q
	return r
}

func (r *Request) QA(args map[string]string) *Request {
	q := url.Values{}
	for k, v := range args {
		q.Add(k, v)
	}

	r.Path += q.Encode()
	return r
}

func (r *Request) C(cert string) *Request {
	r.Certificate = cert
	return r
}

func (r *Request) SC() *Request {
	return r.C(sys_cert_pool)
}

func (r *Request) S(h string, key []byte) *Request {
	r.sign = key
	r.signH = h
	return r
}

func (r *Request) M(m string) *Request {
	r.Method = m
	return r
}

func (r *Request) H(k, v string) *Request {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[k] = v
	return r
}

func (r *Request) B(b interface{}) *Request {
	r.Body = b
	return r
}

func (r *Request) OK(status int) *Request {
	r.Status = status
	return r
}

func (rq *Request) String() string {
	return fmt.Sprintf("%s:%s", rq.Method, rq.URL())
}

func (rq *Request) Hdrs() string {
	ret := ""
	for h, v := range rq.Headers {
		ret += fmt.Sprintf("%s=%s;", h, v)
	}

	return ret
}

func (rq *Request) Tmo(t time.Duration) *Request {
	rq.Timeout = t
	return rq
}

func (rq *Request) make_cert_pool() (*x509.CertPool, error) {
	if rq.Certificate != sys_cert_pool {
		cert_pool := x509.NewCertPool()
		cert_pool.AppendCertsFromPEM([]byte(rq.Certificate))
		return cert_pool, nil
	}

	if runtime.GOOS != "windows" {
		return x509.SystemCertPool()
	} else {
		return x509.NewCertPool(), nil
	}
}

func (rq *Request) make_client() (*http.Client, error) {
	if rq.Certificate == "" {
		return &http.Client{}, nil
	}

	cert_pool, err := rq.make_cert_pool()
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: cert_pool,
			},
		},
	}, nil
}

func (rq *Request) Do() *Response {
	client, err := rq.make_client()
	if err != nil {
		return &Response{err: fmt.Errorf("Cannot make http client: %s", err.Error()) }
	}

	if rq.Timeout != 0 {
		client.Timeout = rq.Timeout
	}

	var body io.Reader

	if rq.Body != nil {
		data, err := json.Marshal(rq.Body)
		if err != nil {
			return &Response{err: fmt.Errorf("Cannot marshal body: %s", err.Error())}
		}

		body = bytes.NewBuffer(data)
		rq.H("Content-Type", "application/json; charset=utf-8")
	}

	http_rq, err := http.NewRequest(rq.Method, rq.URL(), body)
	if err != nil {
		return &Response{err: fmt.Errorf("Cannot make request: %s", err.Error())}
	}

	if rq.sign != nil {
		http_rq.Header.Set(rq.signH, util.HashHMAC(rq.sign, []byte(rq.Method), []byte(rq.Path)))
	}

	for k, v := range rq.Headers {
		http_rq.Header.Set(k, v)
	}

	resp, err := client.Do(http_rq)
	if err != nil {
		return &Response{err: fmt.Errorf("Cannot do http: %s", err.Error())}
	}
	if resp.StatusCode != rq.Status {
		err = fmt.Errorf("Unexpected response: %d", resp.StatusCode)
	}

	return &Response{resp: resp, err: err}
}
