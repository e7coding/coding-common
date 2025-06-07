// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dcrc32 提供 CRC32 校验算法的实用接口
package jcrc32

import (
	"hash/crc32"

	"github.com/e7coding/coding-common/jutil/jconv"
)

// Enc 使用 CRC32 算法对任意类型的变量进行校验计算。
// 会通过 gconv 包将 `v` 转换为字节切片。
func Enc(v interface{}) uint32 {
	return crc32.ChecksumIEEE(jconv.Bytes(v))
}
