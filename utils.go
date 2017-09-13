package main

import (
	"os"
	"os/exec"
	"strings"
	"strconv"
	"log"
)

type Error struct {
	err       error
	CmdOutput string
}

// error string
func (e Error) Error() string {
	return e.err.Error()
}

// make error type
func MakeErr(err error, out string) *Error {
	if err != nil {
		return &Error{
			err:       err,
			CmdOutput: out,
		}
	}
	return nil
}

// to determine whether a file exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// find dbkit backup process id
func findPid() int {
	pidByte, err := exec.Command("/bin/sh", "-c", `ps -eaf|grep "backup start"|grep -v "grep"|awk '{print $2}'`).Output()
	if err != nil {
		log.Fatal(err)
		return -1
	}
	if len(pidByte) == 0 {
		return -1
	}
	pid := string(pidByte)
	pid = strings.TrimSuffix(string(pidByte), "\n")
	if len(pid) == 0 {
		return -1
	}
	intVal, _ := strconv.Atoi(pid)
	return intVal
}
