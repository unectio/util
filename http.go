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
