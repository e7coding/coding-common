// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// ELFHash32 对输入字节切片执行经典 ELF 哈希算法，返回 32 位哈希值。
func ELFHash32(data []byte) uint32 {
	var (
		hash uint32
		x    uint32
	)
	for _, b := range data {
		// 左移 4 位并加上当前字节
		hash = (hash << 4) + uint32(b)
		// 取高 4 位
		if x = hash & 0xF0000000; x != 0 {
			// 将高 4 位右移 24 位后与 hash 异或
			hash ^= x >> 24
			// 清除这些高 4 位
			hash &= ^x
		}
	}
	return hash
}

// ELFHash64 对输入字节切片执行经典 ELF 哈希算法，返回 64 位哈希值。
func ELFHash64(data []byte) uint64 {
	var (
		hash uint64
		x    uint64
	)
	for _, b := range data {
		// 左移 4 位并加上当前字节
		hash = (hash << 4) + uint64(b)
		// 取高 4 位（对应 0xF000...0000）
		if x = hash & 0xF000000000000000; x != 0 {
			// 将高 4 位右移 24 位后与 hash 异或
			hash ^= x >> 24
			// 清除这些高 4 位
			hash &= ^x
		}
	}
	return hash
}
