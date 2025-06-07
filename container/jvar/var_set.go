// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import (
	"github.com/e7coding/coding-common/container/jatomic"
)

// Set sets `value` to `v`, and returns the old value.
func (v *Var) Set(value interface{}) (old interface{}) {
	if v.safe {
		if t, ok := v.val.(*jatomic.Interface); ok {
			old = t.Store(value)
			return
		}
	}
	old = v.val
	v.val = value
	return
}
