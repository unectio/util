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
