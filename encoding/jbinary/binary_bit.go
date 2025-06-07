// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jbinary

// Bit 表示二进制位 (0 或 1)
// 实验性功能，请谨慎使用！
type Bit int8

// EncodeBits 将整数 i 按默认长度 l 编码为比特序列，并可选追加到已有 bits 切片末尾。
func EncodeBits(bits []Bit, i int, l int) []Bit {
	return EncodeBitsWithUint(bits, uint(i), l)
}

// EncodeBitsWithUint 将无符号整数 ui 按指定长度 l 编码为比特序列，
// 并可选追加到已有 bits 切片末尾。
func EncodeBitsWithUint(bits []Bit, ui uint, l int) []Bit {
	a := make([]Bit, l)
	// 从最低位开始填充，直到最高位
	for idx := l - 1; idx >= 0; idx-- {
		a[idx] = Bit(ui & 1)
		ui >>= 1
	}
	if bits != nil {
		return append(bits, a...)
	}
	return a
}

// EncodeBitsToBytes 将比特序列 bits 编码为字节切片，
// 如果 bits 长度不是 8 的倍数，则在末尾补 0。
func EncodeBitsToBytes(bits []Bit) []byte {
	// 补全到字节边界
	if mod := len(bits) % 8; mod != 0 {
		for i := 0; i < 8-mod; i++ {
			bits = append(bits, 0)
		}
	}
	var b []byte
	// 每 8 位为一字节，按大端顺序解析
	for i := 0; i < len(bits); i += 8 {
		b = append(b, byte(DecodeBitsToUint(bits[i:i+8])))
	}
	return b
}

// DecodeBits 将比特序列 bits 解码为 int 值。
func DecodeBits(bits []Bit) int {
	v := 0
	for _, bit := range bits {
		v = (v << 1) | int(bit)
	}
	return v
}

// DecodeBitsToUint 将比特序列 bits 解码为 uint 值。
func DecodeBitsToUint(bits []Bit) uint {
	v := uint(0)
	for _, bit := range bits {
		v = (v << 1) | uint(bit)
	}
	return v
}

// DecodeBytesToBits 将字节切片 bs 解析为比特序列。
func DecodeBytesToBits(bs []byte) []Bit {
	var bits []Bit
	for _, b := range bs {
		bits = EncodeBitsWithUint(bits, uint(b), 8)
	}
	return bits
}
