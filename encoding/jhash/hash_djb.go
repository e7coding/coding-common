// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// DJBHash32 对输入字节切片执行经典的 32 位 DJB 哈希算法，返回 32 位哈希值。
func DJBHash32(data []byte) uint32 {
	var hash uint32 = 5381
	for _, b := range data {
		// hash = hash*33 + b
		hash = (hash << 5) + hash + uint32(b)
	}
	return hash
}

// DJBHash64 对输入字节切片执行经典的 64 位 DJB 哈希算法，返回 64 位哈希值。
func DJBHash64(data []byte) uint64 {
	var hash uint64 = 5381
	for _, b := range data {
		// hash = hash*33 + b
		hash = (hash << 5) + hash + uint64(b)
	}
	return hash
}
