// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import (
	"github.com/e7coding/coding-common/internal/utils"
)

// IsNil 是否为 nil
func (v *Var) IsNil() bool { return utils.IsNil(v.Val()) }

// IsEmpty 是否为空
func (v *Var) IsEmpty() bool { return utils.IsEmpty(v.Val()) }

// IsInt 是否为 int 类型
func (v *Var) IsInt() bool { return utils.IsInt(v.Val()) }

// IsUint 是否为 uint 类型
func (v *Var) IsUint() bool { return utils.IsUint(v.Val()) }

// IsFloat 是否为 float 类型
func (v *Var) IsFloat() bool { return utils.IsFloat(v.Val()) }

// IsSlice 是否为 slice 类型
func (v *Var) IsSlice() bool { return utils.IsSlice(v.Val()) }

// IsMap 是否为 map 类型
func (v *Var) IsMap() bool { return utils.IsMap(v.Val()) }

// IsStruct 是否为 struct 类型
func (v *Var) IsStruct() bool { return utils.IsStruct(v.Val()) }
