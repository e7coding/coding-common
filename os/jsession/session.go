// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsession implements manager and storage features for sessions.
package jsession

import (
	"github.com/e7coding/coding-common/errs/jerr"
	"github.com/e7coding/coding-common/jutil/juid"
)

var (
	// ErrorDisabled is used for marking certain interface function not used.
	ErrorDisabled = jerr.WithMsg("this feature is disabled in this storage")
)

// NewSessionId creates and returns a new and unique session id string,
// which is in 32 bytes.
func NewSessionId() string {
	return juid.S()
}
