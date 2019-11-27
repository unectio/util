package util

import (
	"os"
	"sync"
	"syscall"
	"strings"
	"io/ioutil"
	"path/filepath"
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
	return !strings.Contains("/" + path + "/", "/../")
}

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
