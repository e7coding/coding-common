// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// RSHash32 使用经典 RS 哈希算法生成 32 位哈希值。
func RSHash32(data []byte) uint32 {
	var (
		b    uint32 = 378551
		a    uint32 = 63689
		hash uint32 = 0
	)
	for _, c := range data {
		hash = hash*a + uint32(c)
		a *= b
	}
	return hash
}

// RSHash64 使用经典 RS 哈希算法生成 64 位哈希值。
func RSHash64(data []byte) uint64 {
	var (
		b    uint64 = 378551
		a    uint64 = 63689
		hash uint64 = 0
	)
	for _, c := range data {
		hash = hash*a + uint64(c)
		a *= b
	}
	return hash
}
