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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"time"
)

const (
	chunkSize    = 1024
	preallocSize = 4096
)

var ErrTimeout = errors.New("IO Timeout")

type Socket struct {
	f       *os.File
	peer_fd int
}

func id2fd(id string) (int, error) { return strconv.Atoi(id) }
func fd2id(fd int) string          { return strconv.Itoa(fd) }

func makeSocket(fd, peer int) *Socket {
	return &Socket{f: os.NewFile(uintptr(fd), "seqsk"), peer_fd: peer}
}

func Make() (*Socket, error) {
	pair, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_SEQPACKET, 0)
	if err != nil {
		return nil, fmt.Errorf("socketpair: %s", err.Error())
	}

	syscall.CloseOnExec(pair[1])

	return makeSocket(pair[1], pair[0]), nil
}

func Open(id string) (*Socket, error) {
	fd, err := id2fd(id)
	if err != nil {
		return nil, fmt.Errorf("bad ID: %s", err.Error())
	}

	return makeSocket(fd, -1), nil
}

func (sk *Socket) ClosePeer() {
	if sk.peer_fd != -1 {
		syscall.Close(sk.peer_fd)
		sk.peer_fd = -1
	}
}

func (sk *Socket) PFd() string {
	return fd2id(sk.peer_fd)
}

func (sk *Socket) Close() {
	sk.ClosePeer()
	sk.f.Close()
}

func (sk *Socket) SetTimeout(tmo time.Duration) error {
	tv := syscall.NsecToTimeval(tmo.Nanoseconds())
	return syscall.SetsockoptTimeval(int(sk.f.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
}

/*
 * We use seqpacket socket to send small requests in
 * one send+recv call pair (we don't need to read size
 * first). Thus, for the receiver to know when the large
 * request ends we make the data not chunkSize aliquot
 */

func (sk *Socket) RecvRaw() ([]byte, error) {
	ch := make([]byte, chunkSize)
	data := make([]byte, 0, preallocSize)

	for {
		sz, err := sk.f.Read(ch)
		if err != nil {
			if err == io.EOF {
				return nil, err
			} else if sockTimeout(err) {
				return nil, ErrTimeout
			} else {
				return nil, fmt.Errorf("receive error: %s", err.Error())
			}
		}

		if sz == chunkSize {
			data = append(data, ch...)
			continue
		}

		if ch[sz-1] == byte(0) {
			sz -= 1
		}

		return append(data, ch[:sz]...), nil
	}
}

func (sk *Socket) SendRaw(data []byte) error {
	if len(data)%chunkSize == 0 {
		data = append(data, byte(0))
	}

	_, err := writev(sk.f, split(data, chunkSize))
	if err != nil {
		return fmt.Errorf("send error: %s", err.Error())
	}

	return nil
}

func (sk *Socket) Recv(o interface{}) error {
	data, err := sk.RecvRaw()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, o)
	if err != nil {
		return fmt.Errorf("unmarshal error: %s", err.Error())
	}

	return nil
}

func (sk *Socket) Send(o interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err.Error())
	}

	return sk.SendRaw(b)
}
