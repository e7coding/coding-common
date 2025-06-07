// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import "github.com/e7coding/coding-common/jutil/jconv"

// Ints 转换为 []int
func (v *Var) Ints() []int { return jconv.Ints(v.Val()) }

// Int64s 转换为 []int64
func (v *Var) Int64s() []int64 { return jconv.Int64s(v.Val()) }

// Uints 转换为 []uint
func (v *Var) Uints() []uint { return jconv.Uints(v.Val()) }

// Uint64s 转换为 []uint64
func (v *Var) Uint64s() []uint64 { return jconv.Uint64s(v.Val()) }

// Floats 别名 Float64s，转换为 []float64
func (v *Var) Floats() []float64 { return jconv.Floats(v.Val()) }

// Float32s 转换为 []float32
func (v *Var) Float32s() []float32 { return jconv.Float32s(v.Val()) }

// Float64s 转换为 []float64
func (v *Var) Float64s() []float64 { return jconv.Float64s(v.Val()) }

// Strings 转换为 []string
func (v *Var) Strings() []string { return jconv.Strings(v.Val()) }

// Interfaces 转换为 []interface{}
func (v *Var) Interfaces() []interface{} { return jconv.Interfaces(v.Val()) }

// Slice 别名 Interfaces
func (v *Var) Slice() []interface{} { return v.Interfaces() }

// Array 别名 Interfaces
func (v *Var) Array() []interface{} { return v.Interfaces() }

// Vars 转换为 []*Var
func (v *Var) Vars() []*Var {
	arr := jconv.Interfaces(v.Val())
	if len(arr) == 0 {
		return nil
	}
	vars := make([]*Var, len(arr))
	for i, e := range arr {
		vars[i] = New(e, v.safe)
	}
	return vars
}
