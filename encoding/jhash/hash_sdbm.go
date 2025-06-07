// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// SDBM32 使用经典 SDBM 哈希算法生成 32 位哈希值。
func SDBM32(data []byte) uint32 {
	var hash uint32
	for _, b := range data {
		// 等价于: hash = 65599*hash + uint32(b)
		hash = uint32(b) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}

// SDBM64 使用经典 SDBM 哈希算法生成 64 位哈希值。
func SDBM64(data []byte) uint64 {
	var hash uint64
	for _, b := range data {
		// 等价于: hash = 65599*hash + uint64(b)
		hash = uint64(b) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}
