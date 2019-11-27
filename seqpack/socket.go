package seqpack

import (
	"os"
	"io"
	"fmt"
	"time"
	"errors"
	"strconv"
	"syscall"
	"encoding/json"
)

const (
	chunkSize	= 1024
	preallocSize	= 4096
)

var ErrTimeout = errors.New("IO Timeout")

type Socket struct {
	f	*os.File
	peer_fd	int
}

func id2fd(id string) (int, error) { return strconv.Atoi(id) }
func fd2id(fd int) string { return strconv.Itoa(fd) }

func makeSocket(fd, peer int) (*Socket) {
	return &Socket{ f: os.NewFile(uintptr(fd), "seqsk"), peer_fd: peer }
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

func (sk *Socket)ClosePeer() {
	if sk.peer_fd != -1 {
		syscall.Close(sk.peer_fd)
		sk.peer_fd = -1
	}
}

func (sk *Socket)PFd() string {
	return fd2id(sk.peer_fd)
}

func (sk *Socket)Close() {
	sk.ClosePeer()
	sk.f.Close()
}

func (sk *Socket)SetTimeout(tmo time.Duration) error {
	tv := syscall.NsecToTimeval(tmo.Nanoseconds())
	return syscall.SetsockoptTimeval(int(sk.f.Fd()), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
}

/*
 * We use seqpacket socket to send small requests in
 * one send+recv call pair (we don't need to read size
 * first). Thus, for the receiver to know when the large
 * request ends we make the data not chunkSize aliquot
 */

func (sk *Socket)RecvRaw() ([]byte, error) {
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

func (sk *Socket)SendRaw(data []byte) error {
	if len(data) % chunkSize == 0 {
		data = append(data, byte(0))
	}

	_, err := writev(sk.f, split(data, chunkSize))
	if err != nil {
		return fmt.Errorf("send error: %s", err.Error())
	}

	return nil
}

func (sk *Socket)Recv(o interface{}) error {
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

func (sk *Socket)Send(o interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err.Error())
	}

	return sk.SendRaw(b)
}
