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
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s)) //nolint:errcheck
	return hex.EncodeToString(hash.Sum(nil))
}

var rsRunes = []rune("_1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomBytes(l int) ([]byte, error) {
	ri := make([]byte, l)

	_, err := rand.Read(ri)
	if err != nil {
		return nil, err
	}

	return ri, nil
}

func RandomString(l int) (string, error) {
	ri := make([]byte, l)

	_, err := rand.Read(ri)
	if err != nil {
		return "", err
	}

	rs := make([]rune, l)
	for i, j := range ri {
		rs[i] = rsRunes[int(j)%len(rsRunes)]
	}

	return string(rs), nil
}

func HashHMAC(sec []byte, data ...[]byte) string {
	hash := hmac.New(sha256.New, sec)
	for _, dat := range data {
		hash.Write(dat) //nolint:errcheck
	}
	return hex.EncodeToString(hash.Sum(nil))
}
