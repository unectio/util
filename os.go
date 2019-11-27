package util

import (
	"fmt"
	"os/exec"
	"syscall"
	"runtime"
)

func CmdExited(err error) (bool, int) {
	eerr, x := err.(*exec.ExitError)
	if !x {
		return false, -1
	}

	status := eerr.Sys().(syscall.WaitStatus).ExitStatus()
	return true, status
}

func Caller() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	} else {
		return "¯\\_(ツ)_/¯"
	}
}
