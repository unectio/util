package util

import (
	"io"
	"os"
	"bufio"
	"errors"
	"os/exec"
)

func ReadCmdLines(cmd *exec.Cmd) ([]string, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New("cannot get out pipe: " + err.Error())
	}

	err = cmd.Start()
	if err != nil {
		return nil, errors.New("cannot start cmd: " + err.Error())
	}

	defer cmd.Wait()

	return ReadIOLines(stdout)
}

func ReadIOLines(from io.Reader) ([]string, error) {
	var lines []string

	scaner := bufio.NewScanner(from)
	for scaner.Scan() {
		lines = append(lines, scaner.Text())
	}

	err := scaner.Err()
	if err != nil {
		return nil, errors.New("cannot read lines: " + err.Error())
	}

	return lines, nil
}

type closerWrap struct {
	io.Reader
}

func MakeReadCloser(r io.Reader) io.ReadCloser {
	return &closerWrap{r}
}

func (closerWrap)Close() error {
	return nil
}

func SaveFile(dir, file string, dirPerm os.FileMode, from io.Reader, filePerm os.FileMode) error {
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return err
	}

	to, err := os.OpenFile(dir + "/" + file, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, filePerm)
	if err != nil {
		return err
	}

	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}
