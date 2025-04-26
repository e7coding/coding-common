// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dvar 通用变量类型，支持并发安全和多种类型转换
package dvar

import (
	"github.com/coding-common/container/atomic"
	"github.com/coding-common/util/dconv"
)

// Val 获取当前值
func (v *Var) Val() interface{} {
	if v == nil {
		return nil
	}
	if v.safe {
		if gi, ok := v.val.(*atomic.Interface); ok {
			return gi.Load()
		}
	}
	return v.val
}

// Bytes 转为 []byte
func (v *Var) Bytes() []byte {
	return dconv.Bytes(v.Val())
}

// String 转为 string
func (v *Var) String() string {
	return dconv.String(v.Val())
}

// Bool 转为 bool
func (v *Var) Bool() bool {
	return dconv.Bool(v.Val())
}

// Int 转为 int
func (v *Var) Int() int {
	return dconv.Int(v.Val())
}

// Int8 转为 int8
func (v *Var) Int8() int8 {
	return dconv.Int8(v.Val())
}

// Int16 转为 int16
func (v *Var) Int16() int16 {
	return dconv.Int16(v.Val())
}

// Int32 转为 int32
func (v *Var) Int32() int32 {
	return dconv.Int32(v.Val())
}

// Int64 转为 int64
func (v *Var) Int64() int64 {
	return dconv.Int64(v.Val())
}

// Uint 转为 uint
func (v *Var) Uint() uint {
	return dconv.Uint(v.Val())
}

// Uint8 转为 uint8
func (v *Var) Uint8() uint8 {
	return dconv.Uint8(v.Val())
}

// Uint16 转为 uint16
func (v *Var) Uint16() uint16 {
	return dconv.Uint16(v.Val())
}

// Uint32 转为 uint32
func (v *Var) Uint32() uint32 {
	return dconv.Uint32(v.Val())
}

// Uint64 转为 uint64
func (v *Var) Uint64() uint64 {
	return dconv.Uint64(v.Val())
}

// Float32 转为 float32
func (v *Var) Float32() float32 {
	return dconv.Float32(v.Val())
}

// Float64 转为 float64
func (v *Var) Float64() float64 {
	return dconv.Float64(v.Val())
}
