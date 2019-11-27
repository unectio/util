package restmux

import (
	"io"
	"context"
	"net/http"
	"encoding/json"
)

func respondJson(ctx context.Context, w http.ResponseWriter, from interface{}) Error {
	return respondJson2(ctx, w, http.StatusOK, from)
}

func respondJson2(ctx context.Context, w http.ResponseWriter, status int, from interface{}) Error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(from)
	if err != nil {
		return &GenError{http.StatusInternalServerError, err.Error()}
	}

	return nil
}

func Respond(ctx context.Context, w http.ResponseWriter, result interface{}) {
	err := respondJson(ctx, w, result)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}

func read(body io.ReadCloser, into interface{}) Error {
	defer body.Close()
	err := json.NewDecoder(body).Decode(into)
	if err != nil {
		return &GenError{http.StatusBadRequest, err.Error()}
	}

	return nil
}
