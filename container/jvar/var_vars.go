// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import (
	"github.com/e7coding/coding-common/jutil/jconv"
)

// Vars is a slice of *Var.
type Vars []*Var

// Strings 转为 []string
func (vs Vars) Strings() []string {
	s := make([]string, len(vs))
	for i, v := range vs {
		s[i] = v.String()
	}
	return s
}

// Interfaces 转为 []interface{}
func (vs Vars) Interfaces() []interface{} {
	s := make([]interface{}, len(vs))
	for i, v := range vs {
		s[i] = v.Val()
	}
	return s
}

// Float32s 转为 []float32
func (vs Vars) Float32s() []float32 {
	s := make([]float32, len(vs))
	for i, v := range vs {
		s[i] = v.Float32()
	}
	return s
}

// Float64s 转为 []float64
func (vs Vars) Float64s() []float64 {
	s := make([]float64, len(vs))
	for i, v := range vs {
		s[i] = v.Float64()
	}
	return s
}

// Ints 转为 []int
func (vs Vars) Ints() []int {
	s := make([]int, len(vs))
	for i, v := range vs {
		s[i] = v.Int()
	}
	return s
}

// Int8s 转为 []int8
func (vs Vars) Int8s() []int8 {
	s := make([]int8, len(vs))
	for i, v := range vs {
		s[i] = v.Int8()
	}
	return s
}

// Int16s 转为 []int16
func (vs Vars) Int16s() []int16 {
	s := make([]int16, len(vs))
	for i, v := range vs {
		s[i] = v.Int16()
	}
	return s
}

// Int32s 转为 []int32
func (vs Vars) Int32s() []int32 {
	s := make([]int32, len(vs))
	for i, v := range vs {
		s[i] = v.Int32()
	}
	return s
}

// Int64s 转为 []int64
func (vs Vars) Int64s() []int64 {
	s := make([]int64, len(vs))
	for i, v := range vs {
		s[i] = v.Int64()
	}
	return s
}

// Uints 转为 []uint
func (vs Vars) Uints() []uint {
	s := make([]uint, len(vs))
	for i, v := range vs {
		s[i] = v.Uint()
	}
	return s
}

// Uint8s 转为 []uint8
func (vs Vars) Uint8s() []uint8 {
	s := make([]uint8, len(vs))
	for i, v := range vs {
		s[i] = v.Uint8()
	}
	return s
}

// Uint16s 转为 []uint16
func (vs Vars) Uint16s() []uint16 {
	s := make([]uint16, len(vs))
	for i, v := range vs {
		s[i] = v.Uint16()
	}
	return s
}

// Uint32s 转为 []uint32
func (vs Vars) Uint32s() []uint32 {
	s := make([]uint32, len(vs))
	for i, v := range vs {
		s[i] = v.Uint32()
	}
	return s
}

// Uint64s 转为 []uint64
func (vs Vars) Uint64s() []uint64 {
	s := make([]uint64, len(vs))
	for i, v := range vs {
		s[i] = v.Uint64()
	}
	return s
}

// Scan 转换为指定 struct 切片
func (vs Vars) Scan(pointer interface{}, mapping ...map[string]string) error {
	return jconv.Structs(vs.Interfaces(), pointer, mapping...)
}
