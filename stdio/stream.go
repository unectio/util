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

package stdio

import (
	"strconv"
	"syscall"
)

type Stream struct {
	me   int
	peer int
}

func (s *Stream) Read() string {
	var ret string

	b := make([]byte, 1024)
	for {
		sz, _ := syscall.Read(s.me, b)
		if sz <= 0 {
			break
		}
		ret += string(b[:sz])
	}

	return ret
}

func Make(name string) (*Stream, error) {
	fds := make([]int, 2)
	err := syscall.Pipe(fds)
	if err != nil {
		return nil, err
	}

	syscall.CloseOnExec(fds[0])
	err = syscall.SetNonblock(fds[0], true)
	if err != nil {
		syscall.Close(fds[0])
		syscall.Close(fds[1])
		return nil, err
	}

	ret := Stream{}
	ret.me = fds[0]
	ret.peer = fds[1]

	return &ret, nil
}

func (s *Stream) PFd() string {
	return strconv.Itoa(s.peer)
}

func (s *Stream) Close() {
	syscall.Close(s.me)
	syscall.Close(s.peer)
}

func (s *Stream) Fd() uintptr {
	return uintptr(s.me)
}

func (s *Stream) Pd() int {
	return s.peer
}
