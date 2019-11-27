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

