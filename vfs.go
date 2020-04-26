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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RmOk(err error) error {
	if os.IsNotExist(err) {
		err = nil
	}
	return err
}

func RmdirAsyncPrepare(dir string, sub string) (string, error) {
	path := dir + "/" + sub

	ok, err := Exists(path)
	if !ok {
		return "", err
	}

	tpath, err := ioutil.TempDir(dir, ".trash")
	if err != nil {
		return "", err
	}

	rpath := tpath + "/" + strings.Replace(sub, "/", "_", -1)
	err = os.Rename(path, rpath)
	if err != nil {
		return "", err
	}

	return tpath, nil
}

func rmdirAsyncComplete(tpath string, wg *sync.WaitGroup) {
	if tpath != "" {
		go func() {
			os.RemoveAll(tpath)
			if wg != nil {
				wg.Done()
			}
		}()
	}
}

func RmdirAsyncComplete(tpath string) {
	rmdirAsyncComplete(tpath, nil)
}

func rmdirAsync(dir string, sub string, wg *sync.WaitGroup) error {
	tp, err := RmdirAsyncPrepare(dir, sub)
	if err == nil {
		rmdirAsyncComplete(tp, wg)
	}
	return err
}

func RmdirAsync(dir string, sub string) error {
	return rmdirAsync(dir, sub, nil)
}

func Rmdir(dir string, sub string) error {
	var wg sync.WaitGroup
	wg.Add(1)
	err := rmdirAsync(dir, sub, &wg)
	if err == nil {
		wg.Wait()
	}
	return err
}

type DEntry interface {
	Path() string
}

func WalkTree(prefix string, root DEntry, fn func(os.FileInfo, DEntry) DEntry) error {
	dirs := []DEntry{root}

	for len(dirs) != 0 {
		cur := dirs[0]
		dirs = dirs[1:]

		ents, err := ioutil.ReadDir(prefix + "/" + cur.Path())
		if err != nil {
			return err
		}

		for _, ent := range ents {
			e := fn(ent, cur)
			if e != nil && ent.IsDir() {
				dirs = append(dirs, e)
			}
		}
	}

	return nil
}

func StraightPath(path string) bool {
	return !strings.Contains("/"+path+"/", "/../")
}

/* TBD: move out of Windows-compilable stuff

func DU(dir string) (uint64, error) {
        var bytes uint64

        err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
                if err == nil && path != dir {
			st, _ := info.Sys().(*syscall.Stat_t)
			bytes += uint64(st.Blocks << 9)
		}
                return err
        })

        return bytes, err
}
*/

func open_with_real_path(path string) (*os.File, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}

	real_path, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		f.Close()
		return nil, "", err
	}

	return f, real_path, nil
}

func OpenSafe(dir, name string) (*os.File, error) {
	f, r_path, err := open_with_real_path(dir)
	if err != nil {
		return nil, err
	}

	f.Close()

	f, f_path, err := open_with_real_path(dir + "/" + name)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(f_path, r_path) {
		f.Close()
		return nil, errors.New("target file is not in its directory")
	}

	return f, nil
}
