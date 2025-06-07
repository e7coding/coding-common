// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jbinary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"math"

	"github.com/e7coding/coding-common/internal/intlog"
)

// LeEncode 使用小端格式(LittleEndian)将一个或多个值编码为字节切片。
// 支持常见的整数、浮点数、布尔、字符串及字节切片类型，
// 对于其他类型则调用 binary.Write 进行编码，失败时写入其字符串表示。
func LeEncode(values ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, v := range values {
		if v == nil {
			return buf.Bytes()
		}
		switch value := v.(type) {
		case int:
			buf.Write(LeEncodeInt(value))
		case int8:
			buf.Write(LeEncodeInt8(value))
		case int16:
			buf.Write(LeEncodeInt16(value))
		case int32:
			buf.Write(LeEncodeInt32(value))
		case int64:
			buf.Write(LeEncodeInt64(value))
		case uint:
			buf.Write(LeEncodeUint(value))
		case uint8:
			buf.Write(LeEncodeUint8(value))
		case uint16:
			buf.Write(LeEncodeUint16(value))
		case uint32:
			buf.Write(LeEncodeUint32(value))
		case uint64:
			buf.Write(LeEncodeUint64(value))
		case bool:
			buf.Write(LeEncodeBool(value))
		case string:
			buf.Write(LeEncodeString(value))
		case []byte:
			buf.Write(value)
		case float32:
			buf.Write(LeEncodeFloat32(value))
		case float64:
			buf.Write(LeEncodeFloat64(value))
		default:
			// 尝试使用 binary.Write 编码，失败时记录日志并写入字符串表示
			if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
				intlog.Errorf("%+v", err)
				buf.Write(LeEncodeString(fmt.Sprintf("%v", value)))
			}
		}
	}
	return buf.Bytes()
}

// LeEncodeByLength 使用小端格式编码后填充或截断到指定长度。
func LeEncodeByLength(length int, values ...interface{}) []byte {
	b := LeEncode(values...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[:length]
	}
	return b
}

// LeDecode 使用小端格式将字节切片解码到对应的值，values 应传入指针。
func LeDecode(b []byte, values ...interface{}) error {
	buf := bytes.NewBuffer(b)
	for _, v := range values {
		if err := binary.Read(buf, binary.LittleEndian, v); err != nil {
			return jerr.WithMsgErrF(err, "binary.Read failed")
		}
	}
	return nil
}

// LeEncodeString 将字符串编码为字节切片。
func LeEncodeString(s string) []byte {
	return []byte(s)
}

// LeDecodeToString 将字节切片解码为字符串。
func LeDecodeToString(b []byte) string {
	return string(b)
}

// LeEncodeBool 将布尔值编码为单字节：true->1, false->0。
func LeEncodeBool(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// LeEncodeInt 根据数值范围选择合适的字节长度进行编码。
func LeEncodeInt(i int) []byte {
	if i <= math.MaxInt8 {
		return LeEncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return LeEncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return LeEncodeInt32(int32(i))
	}
	return LeEncodeInt64(int64(i))
}

// LeEncodeUint 根据数值范围选择合适的字节长度进行编码。
func LeEncodeUint(i uint) []byte {
	if i <= math.MaxUint8 {
		return LeEncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return LeEncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return LeEncodeUint32(uint32(i))
	}
	return LeEncodeUint64(uint64(i))
}

// LeEncodeInt8 将 int8 编码为单字节。
func LeEncodeInt8(i int8) []byte {
	return []byte{byte(i)}
}

// LeEncodeUint8 将 uint8 编码为单字节。
func LeEncodeUint8(i uint8) []byte {
	return []byte{i}
}

// LeEncodeInt16 将 int16 编码为 2 字节小端格式。
func LeEncodeInt16(i int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}

// LeEncodeUint16 将 uint16 编码为 2 字节小端格式。
func LeEncodeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

// LeEncodeInt32 将 int32 编码为 4 字节小端格式。
func LeEncodeInt32(i int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

// LeEncodeUint32 将 uint32 编码为 4 字节小端格式。
func LeEncodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

// LeEncodeInt64 将 int64 编码为 8 字节小端格式。
func LeEncodeInt64(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// LeEncodeUint64 将 uint64 编码为 8 字节小端格式。
func LeEncodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

// LeEncodeFloat32 将 float32 编码为 4 字节小端格式。
func LeEncodeFloat32(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(f))
	return b
}

// LeEncodeFloat64 将 float64 编码为 8 字节小端格式。
func LeEncodeFloat64(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(f))
	return b
}

// LeDecodeToInt 将字节切片解码为 int，根据长度自动选择对应类型。
func LeDecodeToInt(b []byte) int {
	switch {
	case len(b) < 2:
		return int(LeDecodeToUint8(b))
	case len(b) < 3:
		return int(LeDecodeToUint16(b))
	case len(b) < 5:
		return int(LeDecodeToUint32(b))
	default:
		return int(LeDecodeToUint64(b))
	}
}

// LeDecodeToUint 将字节切片解码为 uint，根据长度自动选择对应类型。
func LeDecodeToUint(b []byte) uint {
	switch {
	case len(b) < 2:
		return uint(LeDecodeToUint8(b))
	case len(b) < 3:
		return uint(LeDecodeToUint16(b))
	case len(b) < 5:
		return uint(LeDecodeToUint32(b))
	default:
		return uint(LeDecodeToUint64(b))
	}
}

// LeDecodeToBool 将字节切片解码为布尔值，全部为 0 视为 false。
func LeDecodeToBool(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	// 如果全部字节为零，则认为 false
	return !bytes.Equal(b, make([]byte, len(b)))
}

// LeDecodeToInt8 将字节切片第一个字节解码为 int8。
func LeDecodeToInt8(b []byte) int8 {
	if len(b) == 0 {
		panic("empty slice given")
	}
	return int8(b[0])
}

// LeDecodeToUint8 将字节切片第一个字节解码为 uint8。
func LeDecodeToUint8(b []byte) uint8 {
	if len(b) == 0 {
		panic("empty slice given")
	}
	return b[0]
}

// LeDecodeToInt16 将字节切片解码为 int16，小端格式，长度不足时自动填充。
func LeDecodeToInt16(b []byte) int16 {
	return int16(binary.LittleEndian.Uint16(LeFillUpSize(b, 2)))
}

// LeDecodeToUint16 将字节切片解码为 uint16，小端格式，长度不足时自动填充。
func LeDecodeToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(LeFillUpSize(b, 2))
}

// LeDecodeToInt32 将字节切片解码为 int32，小端格式，长度不足时自动填充。
func LeDecodeToInt32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

// LeDecodeToUint32 将字节切片解码为 uint32，小端格式，长度不足时自动填充。
func LeDecodeToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(LeFillUpSize(b, 4))
}

// LeDecodeToInt64 将字节切片解码为 int64，小端格式，长度不足时自动填充。
func LeDecodeToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LeDecodeToUint64 将字节切片解码为 uint64，小端格式，长度不足时自动填充。
func LeDecodeToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(LeFillUpSize(b, 8))
}

// LeDecodeToFloat32 将字节切片解码为 float32，小端格式，长度不足时自动填充。
func LeDecodeToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

// LeDecodeToFloat64 将字节切片解码为 float64，小端格式，长度不足时自动填充。
func LeDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LeFillUpSize 将字节切片填充或截断到指定长度，不修改原切片。
func LeFillUpSize(b []byte, l int) []byte {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)
	return c
}
