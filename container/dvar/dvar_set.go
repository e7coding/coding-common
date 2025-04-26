// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar

import (
	"github.com/coding-common/container/atomic"
)

// Set sets `value` to `v`, and returns the old value.
func (v *Var) Set(value interface{}) (old interface{}) {
	if v.safe {
		if t, ok := v.val.(*atomic.Interface); ok {
			old = t.Store(value)
			return
		}
	}
	old = v.val
	v.val = value
	return
}
