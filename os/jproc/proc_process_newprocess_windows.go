// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

//go:build windows

package jproc

import (
	"syscall"

	"github.com/e7coding/coding-common/text/jstr"
)

// Set the underlying parameters directly on the Windows platform
func joinProcessArgs(p *Process) {
	p.SysProcAttr = &syscall.SysProcAttr{}
	p.SysProcAttr.CmdLine = jstr.Join(p.Args, " ")
}
