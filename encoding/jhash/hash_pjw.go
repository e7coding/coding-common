// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jhash

// PJWHash32 对输入字节切片执行经典的 PJW 哈希算法，返回 32 位哈希值。
func PJWHash32(data []byte) uint32 {
	var (
		BitsInUnsignedInt uint32 = 32 // 4 * 8
		ThreeQuarters            = (BitsInUnsignedInt * 3) / 4
		OneEighth                = BitsInUnsignedInt / 8
		HighBits          uint32 = (0xFFFFFFFF) << (BitsInUnsignedInt - OneEighth)
		hash              uint32
		test              uint32
	)
	for i := 0; i < len(data); i++ {
		hash = (hash << OneEighth) + uint32(data[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}

// PJWHash64 对输入字节切片执行经典的 PJW 哈希算法，返回 64 位哈希值。
// 注意：此实现中 BitsInUnsignedInt 仍使用 32 位长度，以保持与原版算法一致。
func PJWHash64(data []byte) uint64 {
	const (
		BitsInUnsignedInt        = 32 // 保持与 PJWHash32 一致
		ThreeQuarters            = (BitsInUnsignedInt * 3) / 4
		OneEighth                = BitsInUnsignedInt / 8
		HighBits          uint64 = 0xFFFFFFFFFFFFFFFF << (BitsInUnsignedInt - OneEighth)
	)
	var (
		hash uint64
		test uint64
	)
	for _, b := range data {
		hash = (hash << OneEighth) + uint64(b)
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}
