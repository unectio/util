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

package seqpack

import (
	"os"
	"syscall"
)

func sockTimeout(err error) bool {
	perr, ok := err.(*os.PathError)
	if !ok {
		return false
	}

	serr, ok := perr.Err.(syscall.Errno)
	if !ok {
		return false
	}

	return serr.Timeout()
}

func writev(f *os.File, data [][]byte) (int, error) {
	written := 0
	for _, b := range data {
		w, err := f.Write(b)
		if err != nil {
			return written + w, err
		}
	}
	return written, nil
}

func split(buf []byte, sz int) [][]byte {
	chunks := make([][]byte, 0, len(buf)/sz+1)

	for len(buf) >= sz {
		var ch []byte
		ch, buf = buf[:sz], buf[sz:]
		chunks = append(chunks, ch)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}

	return chunks
}
