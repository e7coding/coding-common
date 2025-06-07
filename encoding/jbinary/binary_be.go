// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jbinary

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"math"

	"github.com/e7coding/coding-common/internal/intlog"
)

// BeEncode 使用大端序将一个或多个值编码为字节切片。
func BeEncode(values ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, v := range values {
		if v == nil {
			return buf.Bytes()
		}
		switch value := v.(type) {
		case int:
			buf.Write(BeEncodeInt(value))
		case int8:
			buf.Write(BeEncodeInt8(value))
		case int16:
			buf.Write(BeEncodeInt16(value))
		case int32:
			buf.Write(BeEncodeInt32(value))
		case int64:
			buf.Write(BeEncodeInt64(value))
		case uint:
			buf.Write(BeEncodeUint(value))
		case uint8:
			buf.Write(BeEncodeUint8(value))
		case uint16:
			buf.Write(BeEncodeUint16(value))
		case uint32:
			buf.Write(BeEncodeUint32(value))
		case uint64:
			buf.Write(BeEncodeUint64(value))
		case bool:
			buf.Write(BeEncodeBool(value))
		case string:
			buf.Write(BeEncodeString(value))
		case []byte:
			buf.Write(value)
		case float32:
			buf.Write(BeEncodeFloat32(value))
		case float64:
			buf.Write(BeEncodeFloat64(value))
		default:
			// 其他类型通过 fmt.Sprintf 转为字符串后编码
			if err := binary.Write(buf, binary.BigEndian, value); err != nil {
				intlog.Errorf(context.TODO(), "%+v", err)
				buf.Write(BeEncodeString(fmt.Sprintf("%v", value)))
			}
		}
	}
	return buf.Bytes()
}

// BeEncodeByLength 将已编码字节填充或截断到指定长度。
func BeEncodeByLength(length int, values ...interface{}) []byte {
	b := BeEncode(values...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[:length]
	}
	return b
}

// BeDecode 使用大端序将字节切片解码到给定变量。
func BeDecode(b []byte, values ...interface{}) error {
	buf := bytes.NewBuffer(b)
	for _, v := range values {
		if err := binary.Read(buf, binary.BigEndian, v); err != nil {
			return jerr.WithMsgErrF(err, "binary.Read 失败")
		}
	}
	return nil
}

// BeEncodeString 将字符串转换为字节切片。
func BeEncodeString(s string) []byte {
	return []byte(s)
}

// BeDecodeToString 将字节切片转换为字符串。
func BeDecodeToString(b []byte) string {
	return string(b)
}

// BeEncodeBool 将布尔值编码为单字节（1/0）。
func BeEncodeBool(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// BeEncodeInt 根据值范围选择合适宽度的大端整数编码。
func BeEncodeInt(i int) []byte {
	if i <= math.MaxInt8 {
		return BeEncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return BeEncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return BeEncodeInt32(int32(i))
	}
	return BeEncodeInt64(int64(i))
}

// BeEncodeUint 根据值范围选择合适宽度的大端无符号整数编码。
func BeEncodeUint(i uint) []byte {
	if i <= math.MaxUint8 {
		return BeEncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return BeEncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return BeEncodeUint32(uint32(i))
	}
	return BeEncodeUint64(uint64(i))
}

// 以下函数分别对特定类型进行大端编码与解码：
func BeEncodeInt8(i int8) []byte   { return []byte{byte(i)} }
func BeEncodeUint8(i uint8) []byte { return []byte{i} }
func BeEncodeInt16(i int16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(i))
	return b
}
func BeEncodeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}
func BeEncodeInt32(i int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}
func BeEncodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}
func BeEncodeInt64(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}
func BeEncodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}
func BeEncodeFloat32(f float32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, math.Float32bits(f))
	return b
}
func BeEncodeFloat64(f float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	return b
}

// BeDecodeToInt 根据字节长度进行大端整数解码。
func BeDecodeToInt(b []byte) int {
	l := len(b)
	switch {
	case l < 2:
		return int(BeDecodeToUint8(b))
	case l < 3:
		return int(BeDecodeToUint16(b))
	case l < 5:
		return int(BeDecodeToUint32(b))
	}
	return int(BeDecodeToUint64(b))
}

// BeDecodeToUint 根据字节长度进行大端无符号整数解码。
func BeDecodeToUint(b []byte) uint {
	l := len(b)
	switch {
	case l < 2:
		return uint(BeDecodeToUint8(b))
	case l < 3:
		return uint(BeDecodeToUint16(b))
	case l < 5:
		return uint(BeDecodeToUint32(b))
	}
	return uint(BeDecodeToUint64(b))
}

// BeDecodeToBool 将全零字节视为 false，否则为 true。
func BeDecodeToBool(b []byte) bool {
	if len(b) == 0 || bytes.Equal(b, make([]byte, len(b))) {
		return false
	}
	return true
}

func BeDecodeToInt8(b []byte) int8     { return int8(b[0]) }
func BeDecodeToUint8(b []byte) uint8   { return b[0] }
func BeDecodeToInt16(b []byte) int16   { return int16(binary.BigEndian.Uint16(BeFillUpSize(b, 2))) }
func BeDecodeToUint16(b []byte) uint16 { return binary.BigEndian.Uint16(BeFillUpSize(b, 2)) }
func BeDecodeToInt32(b []byte) int32   { return int32(binary.BigEndian.Uint32(BeFillUpSize(b, 4))) }
func BeDecodeToUint32(b []byte) uint32 { return binary.BigEndian.Uint32(BeFillUpSize(b, 4)) }
func BeDecodeToInt64(b []byte) int64   { return int64(binary.BigEndian.Uint64(BeFillUpSize(b, 8))) }
func BeDecodeToUint64(b []byte) uint64 { return binary.BigEndian.Uint64(BeFillUpSize(b, 8)) }
func BeDecodeToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(BeFillUpSize(b, 4)))
}
func BeDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(BeFillUpSize(b, 8)))
}

// BeFillUpSize 将字节填充到指定长度（前置零填充）。
func BeFillUpSize(b []byte, l int) []byte {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c[l-len(b):], b)
	return c
}
