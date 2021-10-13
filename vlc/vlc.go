package vlc

import (
	"os/exec"
	"errors"
	"strings"
	"os"
)


func fileExists (filepath string) bool {
	fileinfo, err := os.Stat(filepath)
    if os.IsNotExist(err) {
		return false
	}
	return !fileinfo.IsDir()
}


func OpenVLC(path string, params string, fileToOpen string) (int, error) {
	pid := -1
	var paramsSlice []string

	if !fileExists(path) {
		return pid, errors.New("Path not exist")
	}

	if params != "" {
		paramsSlice = strings.Split(params, " ")
	}

	paramsSlice = append(paramsSlice, fileToOpen)

	cmd := exec.Command(path, paramsSlice...)
	err := cmd.Start()
	if err != nil {
		return pid, errors.New("Could not open an vlc instance")
	}

	pid = cmd.Process.Pid

	return pid, nil
}

func CloseVLC(pid int) {
	
}
