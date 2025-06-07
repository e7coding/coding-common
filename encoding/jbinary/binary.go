// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dbinary 提供二进制与字节数据处理的便捷接口。
// 默认使用小端（LittleEndian）进行编码和解码。

package jbinary

// Encode 将任意类型的值按照小端编码并返回字节切片。
func Encode(values ...interface{}) []byte {
	return LeEncode(values...)
}

// EncodeByLength 按指定长度进行小端编码并返回字节切片。
func EncodeByLength(length int, values ...interface{}) []byte {
	return LeEncodeByLength(length, values...)
}

// Decode 使用小端解码字节切片到对应的值，values 必须是指针。
func Decode(b []byte, values ...interface{}) error {
	return LeDecode(b, values...)
}

// EncodeString 将字符串按照小端编码为字节切片。
func EncodeString(s string) []byte {
	return LeEncodeString(s)
}

// DecodeToString 将字节切片按照小端解码为字符串。
func DecodeToString(b []byte) string {
	return LeDecodeToString(b)
}

// EncodeBool 将布尔值按照小端编码为字节切片。
func EncodeBool(b bool) []byte {
	return LeEncodeBool(b)
}

// EncodeInt 将 int 类型按照小端编码为字节切片。
func EncodeInt(i int) []byte {
	return LeEncodeInt(i)
}

// EncodeUint 将 uint 类型按照小端编码为字节切片。
func EncodeUint(i uint) []byte {
	return LeEncodeUint(i)
}

// EncodeInt8 将 int8 类型按照小端编码为字节切片。
func EncodeInt8(i int8) []byte {
	return LeEncodeInt8(i)
}

// EncodeUint8 将 uint8 类型按照小端编码为字节切片。
func EncodeUint8(i uint8) []byte {
	return LeEncodeUint8(i)
}

// EncodeInt16 将 int16 类型按照小端编码为字节切片。
func EncodeInt16(i int16) []byte {
	return LeEncodeInt16(i)
}

// EncodeUint16 将 uint16 类型按照小端编码为字节切片。
func EncodeUint16(i uint16) []byte {
	return LeEncodeUint16(i)
}

// EncodeInt32 将 int32 类型按照小端编码为字节切片。
func EncodeInt32(i int32) []byte {
	return LeEncodeInt32(i)
}

// EncodeUint32 将 uint32 类型按照小端编码为字节切片。
func EncodeUint32(i uint32) []byte {
	return LeEncodeUint32(i)
}

// EncodeInt64 将 int64 类型按照小端编码为字节切片。
func EncodeInt64(i int64) []byte {
	return LeEncodeInt64(i)
}

// EncodeUint64 将 uint64 类型按照小端编码为字节切片。
func EncodeUint64(i uint64) []byte {
	return LeEncodeUint64(i)
}

// EncodeFloat32 将 float32 类型按照小端编码为字节切片。
func EncodeFloat32(f float32) []byte {
	return LeEncodeFloat32(f)
}

// EncodeFloat64 将 float64 类型按照小端编码为字节切片。
func EncodeFloat64(f float64) []byte {
	return LeEncodeFloat64(f)
}

// DecodeToInt 将字节切片按照小端解码为 int 类型。
func DecodeToInt(b []byte) int {
	return LeDecodeToInt(b)
}

// DecodeToUint 将字节切片按照小端解码为 uint 类型。
func DecodeToUint(b []byte) uint {
	return LeDecodeToUint(b)
}

// DecodeToBool 将字节切片按照小端解码为 bool 类型。
func DecodeToBool(b []byte) bool {
	return LeDecodeToBool(b)
}

// DecodeToInt8 将字节切片按照小端解码为 int8 类型。
func DecodeToInt8(b []byte) int8 {
	return LeDecodeToInt8(b)
}

// DecodeToUint8 将字节切片按照小端解码为 uint8 类型。
func DecodeToUint8(b []byte) uint8 {
	return LeDecodeToUint8(b)
}

// DecodeToInt16 将字节切片按照小端解码为 int16 类型。
func DecodeToInt16(b []byte) int16 {
	return LeDecodeToInt16(b)
}

// DecodeToUint16 将字节切片按照小端解码为 uint16 类型。
func DecodeToUint16(b []byte) uint16 {
	return LeDecodeToUint16(b)
}

// DecodeToInt32 将字节切片按照小端解码为 int32 类型。
func DecodeToInt32(b []byte) int32 {
	return LeDecodeToInt32(b)
}

// DecodeToUint32 将字节切片按照小端解码为 uint32 类型。
func DecodeToUint32(b []byte) uint32 {
	return LeDecodeToUint32(b)
}

// DecodeToInt64 将字节切片按照小端解码为 int64 类型。
func DecodeToInt64(b []byte) int64 {
	return LeDecodeToInt64(b)
}

// DecodeToUint64 将字节切片按照小端解码为 uint64 类型。
func DecodeToUint64(b []byte) uint64 {
	return LeDecodeToUint64(b)
}

// DecodeToFloat32 将字节切片按照小端解码为 float32 类型。
func DecodeToFloat32(b []byte) float32 {
	return LeDecodeToFloat32(b)
}

// DecodeToFloat64 将字节切片按照小端解码为 float64 类型。
func DecodeToFloat64(b []byte) float64 {
	return LeDecodeToFloat64(b)
}
