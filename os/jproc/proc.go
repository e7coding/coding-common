// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gproc implements management and communication for processes.
package jproc

import (
	"os"
	"runtime"
	"time"

	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/os/jenv"
	"github.com/e7coding/coding-common/os/jfile"
	"github.com/e7coding/coding-common/text/jstr"
)

const (
	envKeyPPid            = "GPROC_PPID"
	tracingInstrumentName = "github.com/e7coding/coding-common/os/gproc.Process"
)

var (
	processPid       = os.Getpid() // processPid is the pid of current process.
	processStartTime = time.Now()  // processStartTime is the start time of current process.
)

// Pid returns the pid of current process.
func Pid() int {
	return processPid
}

// PPid returns the custom parent pid if exists, or else it returns the system parent pid.
func PPid() int {
	if !IsChild() {
		return Pid()
	}
	ppidValue := os.Getenv(envKeyPPid)
	if ppidValue != "" && ppidValue != "0" {
		return jconv.Int(ppidValue)
	}
	return PPidOS()
}

// PPidOS returns the system parent pid of current process.
// Note that the difference between PPidOS and PPid function is that the PPidOS returns
// the system ppid, but the PPid functions may return the custom pid by gproc if the custom
// ppid exists.
func PPidOS() int {
	return os.Getppid()
}

// IsChild checks and returns whether current process is a child process.
// A child process is forked by another gproc process.
func IsChild() bool {
	ppidValue := os.Getenv(envKeyPPid)
	return ppidValue != "" && ppidValue != "0"
}

// SetPPid sets custom parent pid for current process.
func SetPPid(ppid int) error {
	if ppid > 0 {
		return os.Setenv(envKeyPPid, jconv.String(ppid))
	} else {
		return os.Unsetenv(envKeyPPid)
	}
}

// StartTime returns the start time of current process.
func StartTime() time.Time {
	return processStartTime
}

// Uptime returns the duration which current process has been running
func Uptime() time.Duration {
	return time.Since(processStartTime)
}

// SearchBinary searches the binary `file` in current working folder and PATH environment.
func SearchBinary(file string) string {
	// Check if it is absolute path of exists at current working directory.
	if jfile.Exists(file) {
		return file
	}
	return SearchBinaryPath(file)
}

// SearchBinaryPath searches the binary `file` in PATH environment.
func SearchBinaryPath(file string) string {
	array := ([]string)(nil)
	switch runtime.GOOS {
	case "windows":
		envPath := jenv.Get("PATH", jenv.Get("Path")).String()
		if jstr.Contains(envPath, ";") {
			array = jstr.SplitAndTrim(envPath, ";")
		} else if jstr.Contains(envPath, ":") {
			array = jstr.SplitAndTrim(envPath, ":")
		}
		if jfile.Ext(file) != ".exe" {
			file += ".exe"
		}

	default:
		array = jstr.SplitAndTrim(jenv.Get("PATH").String(), ":")
	}
	if len(array) > 0 {
		path := ""
		for _, v := range array {
			path = v + jfile.Separator + file
			if jfile.Exists(path) && jfile.IsFile(path) {
				return path
			}
		}
	}
	return ""
}
