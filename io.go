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
	"bufio"
	"errors"
	"io"
	"os"
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

	defer cmd.Wait() //nolint:errcheck

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

func (closerWrap) Close() error {
	return nil
}

func SaveFile(dir, file string, dirPerm os.FileMode, from io.Reader, filePerm os.FileMode) error {
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return err
	}

	to, err := os.OpenFile(dir+"/"+file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
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
