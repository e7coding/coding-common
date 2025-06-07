// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dvar 通用变量类型，支持并发安全
package jvar

import (
	"github.com/e7coding/coding-common/container/jatomic"
	"github.com/e7coding/coding-common/internal/json"
)

// Var 通用变量类型
type Var struct {
	val  interface{} // 底层值
	safe bool        // 并发安全标志
}

// New 创建 Var，可选开启并发安全
func New(initial interface{}, safe ...bool) *Var {
	if len(safe) > 0 && safe[0] {
		return &Var{
			val:  jatomic.NewInterface(initial),
			safe: true,
		}
	}
	return &Var{val: initial}
}

// Value 获取当前值
func (v *Var) Value() interface{} {
	if v.safe {
		return v.val.(*jatomic.Interface).Load()
	}
	return v.val
}

// MarshalJSON JSON 序列化
func (v *Var) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value())
}

// UnmarshalJSON JSON 反序列化
func (v *Var) UnmarshalJSON(b []byte) error {
	var tmp interface{}
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	v.Set(tmp)
	return nil
}

// UnmarshalValue 解析并设置任意类型的值
func (v *Var) UnmarshalValue(value interface{}) error {
	v.Set(value)
	return nil
}
