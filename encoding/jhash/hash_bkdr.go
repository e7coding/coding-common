// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// BKDRHash32 对输入字节切片执行经典的 32 位 BKDR 哈希算法，返回 32 位哈希值。
func BKDRHash32(data []byte) uint32 {
	const seed uint32 = 131
	var hash uint32
	for _, b := range data {
		hash = hash*seed + uint32(b)
	}
	return hash
}

// BKDRHash64 对输入字节切片执行经典的 64 位 BKDR 哈希算法，返回 64 位哈希值。
func BKDRHash64(data []byte) uint64 {
	const seed uint64 = 131
	var hash uint64
	for _, b := range data {
		hash = hash*seed + uint64(b)
	}
	return hash
}
