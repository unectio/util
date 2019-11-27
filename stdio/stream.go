package stdio

import (
	"strconv"
	"syscall"
)

type Stream struct {
	me	int
	peer	int
}

func (s *Stream)Read() string {
	var ret string

	b := make([]byte, 1024, 1024)
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

func (s *Stream)PFd() string {
	return strconv.Itoa(s.peer)
}

func (s *Stream)Close() {
	syscall.Close(s.me)
	syscall.Close(s.peer)
}

func (s *Stream)Fd() uintptr {
	return uintptr(s.me)
}

func (s *Stream)Pd() int {
	return s.peer
}
