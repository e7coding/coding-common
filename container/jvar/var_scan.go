// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import (
	"github.com/e7coding/coding-common/jutil/jconv"
)

// Scan 自动识别并填充变量指针，支持 map 结构映射
func (v *Var) Scan(pointer interface{}, mapping ...map[string]string) error {
	return jconv.Scan(v.Val(), pointer, mapping...)
}
