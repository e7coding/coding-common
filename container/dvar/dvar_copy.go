// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar

import (
	"github.com/coding-common/internal/deepcopy"
	"github.com/coding-common/util/dutil"
)

// Clone 浅拷贝 Var
func (v *Var) Clone() *Var {
	return New(v.Val(), v.safe)
}

// Copy 深拷贝 Var（数据级别复制）
func (v *Var) Copy() *Var {
	return New(dutil.Copy(v.Val()), v.safe)
}

// DeepCopy 实现深拷贝接口
func (v *Var) DeepCopy() interface{} {
	if v == nil {
		return nil
	}
	return New(deepcopy.Copy(v.Val()), v.safe)
}
