package rq

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/unectio/util"
)

type Request struct {
	Method		string
	Host		string
	Path		string
	Headers		map[string]string
	Status		int
	Body		interface{}
	Timeout		time.Duration

	sign		[]byte
	signH		string
}

func (rq *Request)URL() string {
	return rq.Host + rq.Path
}

func Req(host, url string) *Request {
	return &Request {
		Host:		host,
		Path:		url,

		Method:		"POST",
		Status:		http.StatusOK,
	}
}

func (r *Request)Q(q string) *Request {
	r.Path += q
	return r
}

func (r *Request)S(h string, key []byte) *Request {
	r.sign = key
	r.signH = h
	return r
}

func (r *Request)M(m string) *Request {
	r.Method = m
	return r
}

func (r *Request)H(k, v string) *Request {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[k] = v
	return r
}

func (r *Request)B(b interface{}) *Request {
	r.Body = b
	return r
}

func (r *Request)OK(status int) *Request {
	r.Status = status
	return r
}

func (rq *Request)String() string {
	return fmt.Sprintf("%s:%s", rq.Method, rq.URL())
}

func (rq *Request)Hdrs() string {
	ret := ""
	for h, v := range rq.Headers {
		ret += fmt.Sprintf("%s=%s;", h, v)
	}

	return ret
}

func (rq *Request)Tmo(t time.Duration) *Request {
	rq.Timeout = t
	return rq
}

func (rq *Request)Do() *Response {
	client := &http.Client{ }

	if rq.Timeout != 0 {
		client.Timeout = rq.Timeout
	}

	var body io.Reader

	if rq.Body != nil {
		data, err := json.Marshal(rq.Body)
		if err != nil {
			return &Response{ err: fmt.Errorf("Cannot marshal body: %s", err.Error()) }
		}

		body = bytes.NewBuffer(data)
		rq.H("Content-Type", "application/json; charset=utf-8")
	}

	http_rq, err := http.NewRequest(rq.Method, rq.URL(), body)
	if err != nil {
		return &Response{ err: fmt.Errorf("Cannot make request: %s", err.Error()) }
	}

	if rq.sign != nil {
		http_rq.Header.Set(rq.signH, util.HashHMAC(rq.sign, []byte(rq.Method), []byte(rq.Path)))
	}

	for k, v := range rq.Headers {
		http_rq.Header.Set(k, v)
	}

	resp, err := client.Do(http_rq)
	if err != nil {
		return &Response{ err: fmt.Errorf("Cannot do http: %s", err.Error()) }
	}

	if resp.StatusCode != rq.Status {
		err = fmt.Errorf("Unexpected response: %d", resp.StatusCode)
	}

	return &Response{ resp: resp, err: err }
}
