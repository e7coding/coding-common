// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// JSHash32 对输入字节切片执行经典 JS 哈希算法，返回 32 位哈希值。
func JSHash32(data []byte) uint32 {
	var hash uint32 = 1315423911
	for _, b := range data {
		// hash = hash ^ ((hash << 5) + b + (hash >> 2))
		hash ^= (hash << 5) + uint32(b) + (hash >> 2)
	}
	return hash
}

// JSHash64 对输入字节切片执行经典 JS 哈希算法，返回 64 位哈希值。
func JSHash64(data []byte) uint64 {
	var hash uint64 = 1315423911
	for _, b := range data {
		// hash = hash ^ ((hash << 5) + b + (hash >> 2))
		hash ^= (hash << 5) + uint64(b) + (hash >> 2)
	}
	return hash
}
