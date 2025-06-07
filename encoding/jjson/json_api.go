// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jjson

import (
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"

	"github.com/e7coding/coding-common/container/jvar"
)

// Interface 返回 Json 对象的底层值。
func (j *Json) Interface() interface{} {
	if j == nil {
		return nil
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	if j.p == nil {
		return nil
	}
	return *(j.p)
}

// Var 以 *jvar.Var 形式返回 Json 值。
func (j *Json) Var() *jvar.Var {
	return jvar.New(j.Interface())
}

// IsNil 判断 Json 对象是否为空。
func (j *Json) IsNil() bool {
	if j == nil {
		return true
	}
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.p == nil || *(j.p) == nil
}

// Get 按路径 pattern 获取值，未找到时返回默认值 def。
func (j *Json) Get(pattern string, def ...interface{}) *jvar.Var {
	if j == nil {
		return nil
	}
	j.mu.RLock()
	defer j.mu.RUnlock()

	if pattern == "" {
		return nil
	}
	result := j.getPointerByPattern(pattern)
	if result != nil {
		return jvar.New(*result)
	}
	if len(def) > 0 {
		return jvar.New(def[0])
	}
	return nil
}

// GetJson 按路径 pattern 获取子 Json 对象。
func (j *Json) GetJson(pattern string, def ...interface{}) *Json {
	return New(j.Get(pattern, def...).Val())
}

// GetJsons 按路径 pattern 获取 Json 切片。
func (j *Json) GetJsons(pattern string, def ...interface{}) []*Json {
	array := j.Get(pattern, def...).Array()
	if len(array) > 0 {
		jsonSlice := make([]*Json, len(array))
		for i := range array {
			jsonSlice[i] = New(array[i])
		}
		return jsonSlice
	}
	return nil
}

// GetJsonMap 按路径 pattern 获取 Json 对象映射。
func (j *Json) GetJsonMap(pattern string, def ...interface{}) map[string]*Json {
	m := j.Get(pattern, def...).Map()
	if len(m) > 0 {
		jsonMap := make(map[string]*Json, len(m))
		for k, v := range m {
			jsonMap[k] = New(v)
		}
		return jsonMap
	}
	return nil
}

// Set 按路径 pattern 设置值 value。
func (j *Json) Set(pattern string, value interface{}) error {
	return j.setValue(pattern, value, false)
}

// MustSet 与 Set 相同，但发生错误时 panic。
func (j *Json) MustSet(pattern string, value interface{}) {
	if err := j.Set(pattern, value); err != nil {
		panic(err)
	}
}

// Remove 按路径 pattern 删除值。
func (j *Json) Remove(pattern string) error {
	return j.setValue(pattern, nil, true)
}

// MustRemove 与 Remove 相同，但发生错误时 panic。
func (j *Json) MustRemove(pattern string) {
	if err := j.Remove(pattern); err != nil {
		panic(err)
	}
}

// Contains 检查路径 pattern 对应的值是否存在。
func (j *Json) Contains(pattern string) bool {
	return j.Get(pattern) != nil
}

// Len 返回路径 pattern 对应值的长度，仅对 slice 或 map 有效，不存在或类型不符返回 -1。
func (j *Json) Len(pattern string) int {
	p := j.getPointerByPattern(pattern)
	if p != nil {
		switch (*p).(type) {
		case map[string]interface{}:
			return len((*p).(map[string]interface{}))
		case []interface{}:
			return len((*p).([]interface{}))
		default:
			return -1
		}
	}
	return -1
}

// Append 向路径 pattern 对应的 slice 追加值 value，若不存在则创建。
func (j *Json) Append(pattern string, value interface{}) error {
	p := j.getPointerByPattern(pattern)
	if p == nil || *p == nil {
		if pattern == "." {
			return j.Set("0", value)
		}
		return j.Set(fmt.Sprintf("%s.0", pattern), value)
	}
	switch (*p).(type) {
	case []interface{}:
		if pattern == "." {
			return j.Set(fmt.Sprintf("%d", len((*p).([]interface{}))), value)
		}
		return j.Set(fmt.Sprintf("%s.%d", pattern, len((*p).([]interface{}))), value)
	}
	return jerr.WithMsgF("无效的变量类型 %s", pattern)
}

// MustAppend 与 Append 相同，但发生错误时 panic。
func (j *Json) MustAppend(pattern string, value interface{}) {
	if err := j.Append(pattern, value); err != nil {
		panic(err)
	}
}

// Map 将 Json 对象转换为 map[string]interface{}。
func (j *Json) Map() map[string]interface{} {
	return j.Var().Map()
}

// Array 将 Json 对象转换为 []interface{}。
func (j *Json) Array() []interface{} {
	return j.Var().Array()
}

// Scan 自动根据 pointer 类型调用 Struct 或 Structs 方法进行转换。
func (j *Json) Scan(pointer interface{}, mapping ...map[string]string) error {
	return j.Var().Scan(pointer, mapping...)
}
