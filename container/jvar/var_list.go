// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import (
	"github.com/e7coding/coding-common/jutil"
)

// ListItemValues 提取 slice 中每个元素（map 或 struct）的指定字段值
func (v *Var) ListItemValues(key interface{}) []interface{} {
	return jutil.ListItemValues(v.Val(), key)
}

// ListItemValuesUnique 提取并去重 slice 中每个元素的指定字段值
func (v *Var) ListItemValuesUnique(key string) []interface{} {
	return jutil.ListItemValuesUnique(v.Val(), key)
}
