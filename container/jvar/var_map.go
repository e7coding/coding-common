// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import "github.com/e7coding/coding-common/jutil/jconv"

// MapOption Map 转换选项别名
type MapOption = jconv.MapOption

// Map 转换为 map[string]interface{}
func (v *Var) Map(opts ...MapOption) map[string]interface{} {
	return jconv.Map(v.Val(), opts...)
}

// MapStrAny 同 Map
func (v *Var) MapStrAny(opts ...MapOption) map[string]interface{} {
	return v.Map(opts...)
}

// MapStrStr 转换为 map[string]string
func (v *Var) MapStrStr(opts ...MapOption) map[string]string {
	return jconv.MapStrStr(v.Val(), opts...)
}

// MapStrVar 转换为 map[string]*Var
func (v *Var) MapStrVar(opts ...MapOption) map[string]*Var {
	m := v.Map(opts...)
	if len(m) == 0 {
		return nil
	}
	vm := make(map[string]*Var, len(m))
	for k, val := range m {
		vm[k] = New(val, v.safe)
	}
	return vm
}

// Maps 转换为 []map[string]interface{}
func (v *Var) Maps(opts ...MapOption) []map[string]interface{} {
	return jconv.Maps(v.Val(), opts...)
}

// MapToMap 转换到指定 map 变量
func (v *Var) MapToMap(pointer interface{}, mapping ...map[string]string) error {
	return jconv.MapToMap(v.Val(), pointer, mapping...)
}

// MapToMaps 转换到指定 []map 变量
func (v *Var) MapToMaps(pointer interface{}, mapping ...map[string]string) error {
	return jconv.MapToMaps(v.Val(), pointer, mapping...)
}

// MapToMapsDeep 递归转换到指定 []map 变量
func (v *Var) MapToMapsDeep(pointer interface{}, mapping ...map[string]string) error {
	return jconv.MapToMap(v.Val(), pointer, mapping...)
}
