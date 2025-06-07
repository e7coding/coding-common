// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jvar

import (
	"time"

	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/os/jtime"
)

// Time 转为 time.Time，格式为 gtime 格式，如 Y-m-d H:i:s
func (v *Var) Time(format ...string) time.Time { return jconv.Time(v.Val(), format...) }

// Duration 转为 time.Duration，字符串使用 time.ParseDuration
func (v *Var) Duration() time.Duration { return jconv.Duration(v.Val()) }

// GTime 转为 *time.Time，格式为 gtime 格式，如 Y-m-d H:i:s
func (v *Var) GTime(format ...string) *jtime.Time { return jconv.GTime(v.Val(), format...) }
