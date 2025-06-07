// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jutil

import (
	"strings"

	"github.com/e7coding/coding-common/jutil/jconv"
)

// Comparator is a function that compare a and b, and returns the result as int.
//
// Should return a number:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Comparator func(a, b interface{}) int

// ComparatorString provides a fast comparison on strings.
func ComparatorString(a, b interface{}) int {
	return strings.Compare(jconv.String(a), jconv.String(b))
}

// ComparatorInt provides a basic comparison on int.
func ComparatorInt(a, b interface{}) int {
	return jconv.Int(a) - jconv.Int(b)
}

// ComparatorInt8 provides a basic comparison on int8.
func ComparatorInt8(a, b interface{}) int {
	return int(jconv.Int8(a) - jconv.Int8(b))
}

// ComparatorInt16 provides a basic comparison on int16.
func ComparatorInt16(a, b interface{}) int {
	return int(jconv.Int16(a) - jconv.Int16(b))
}

// ComparatorInt32 provides a basic comparison on int32.
func ComparatorInt32(a, b interface{}) int {
	return int(jconv.Int32(a) - jconv.Int32(b))
}

// ComparatorInt64 provides a basic comparison on int64.
func ComparatorInt64(a, b interface{}) int {
	return int(jconv.Int64(a) - jconv.Int64(b))
}

// ComparatorUint provides a basic comparison on uint.
func ComparatorUint(a, b interface{}) int {
	return int(jconv.Uint(a) - jconv.Uint(b))
}

// ComparatorUint8 provides a basic comparison on uint8.
func ComparatorUint8(a, b interface{}) int {
	return int(jconv.Uint8(a) - jconv.Uint8(b))
}

// ComparatorUint16 provides a basic comparison on uint16.
func ComparatorUint16(a, b interface{}) int {
	return int(jconv.Uint16(a) - jconv.Uint16(b))
}

// ComparatorUint32 provides a basic comparison on uint32.
func ComparatorUint32(a, b interface{}) int {
	return int(jconv.Uint32(a) - jconv.Uint32(b))
}

// ComparatorUint64 provides a basic comparison on uint64.
func ComparatorUint64(a, b interface{}) int {
	return int(jconv.Uint64(a) - jconv.Uint64(b))
}

// ComparatorFloat32 provides a basic comparison on float32.
func ComparatorFloat32(a, b interface{}) int {
	aFloat := jconv.Float32(a)
	bFloat := jconv.Float32(b)
	if aFloat == bFloat {
		return 0
	}
	if aFloat > bFloat {
		return 1
	}
	return -1
}

// ComparatorFloat64 provides a basic comparison on float64.
func ComparatorFloat64(a, b interface{}) int {
	aFloat := jconv.Float64(a)
	bFloat := jconv.Float64(b)
	if aFloat == bFloat {
		return 0
	}
	if aFloat > bFloat {
		return 1
	}
	return -1
}

// ComparatorByte provides a basic comparison on byte.
func ComparatorByte(a, b interface{}) int {
	return int(jconv.Byte(a) - jconv.Byte(b))
}

// ComparatorRune provides a basic comparison on rune.
func ComparatorRune(a, b interface{}) int {
	return int(jconv.Rune(a) - jconv.Rune(b))
}

// ComparatorTime provides a basic comparison on time.Time.
func ComparatorTime(a, b interface{}) int {
	aTime := jconv.Time(a)
	bTime := jconv.Time(b)
	switch {
	case aTime.After(bTime):
		return 1
	case aTime.Before(bTime):
		return -1
	default:
		return 0
	}
}
