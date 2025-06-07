// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// APHash32 对输入字节切片执行经典的 32 位 AP 哈希算法，返回 32 位哈希值。
func APHash32(data []byte) uint32 {
	var hash uint32
	for i, b := range data {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint32(b) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint32(b) ^ (hash >> 5)) + 1
		}
	}
	return hash
}

// APHash64 对输入字节切片执行经典的 64 位 AP 哈希算法，返回 64 位哈希值。
func APHash64(data []byte) uint64 {
	var hash uint64
	for i, b := range data {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint64(b) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint64(b) ^ (hash >> 5)) + 1
		}
	}
	return hash
}
