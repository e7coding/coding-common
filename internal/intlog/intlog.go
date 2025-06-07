// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package intlog provides internal logging for GoFrame development usage only.
package intlog

import (
	"bytes"
	"fmt"
	"time"

	"github.com/e7coding/coding-common/internal/utils"
)

// Print prints `v` with newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Print(v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(fmt.Sprint(v...))
}

// Printf prints `v` with format `format` using fmt.Printf.
// The parameter `v` can be multiple variables.
func Printf(format string, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(fmt.Sprintf(format, v...))
}

// Error prints `v` with newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Error(v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(fmt.Sprint(v...))
}

// Errorf prints `v` with format `format` using fmt.Printf.
func Errorf(format string, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(fmt.Sprintf(format, v...))
}

// PrintFunc prints the output from function `f`.
// It only calls function `f` if debug mode is enabled.
func PrintFunc(f func() string) {
	if !utils.IsDebugEnabled() {
		return
	}
	s := f()
	if s == "" {
		return
	}
	doPrint(s)
}

// ErrorFunc prints the output from function `f`.
// It only calls function `f` if debug mode is enabled.
func ErrorFunc(f func() string) {
	if !utils.IsDebugEnabled() {
		return
	}
	s := f()
	if s == "" {
		return
	}
	doPrint(s)
}

func doPrint(content string) {

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05.000"))
	buffer.WriteString(" [INTO] ")
	buffer.WriteString(" ")
	buffer.WriteString(content)
	buffer.WriteString("\n")
	fmt.Print(buffer.String())
}
