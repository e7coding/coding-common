// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar

import (
	"github.com/coding-common/util/dconv"
)

// Struct 转换为指定 struct 实例
func (v *Var) Struct(pointer interface{}, mapping ...map[string]string) error {
	return dconv.Struct(v.Val(), pointer, mapping...)
}

// Structs 转换为指定 struct 切片
func (v *Var) Structs(pointer interface{}, mapping ...map[string]string) error {
	return dconv.Structs(v.Val(), pointer, mapping...)
}
