// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jconv

import (
	"time"

	"github.com/e7coding/coding-common/os/jtime"
)

// Time converts `any` to time.Time.
func Time(any any, format ...string) time.Time {
	t, _ := defaultConverter.Time(any, format...)
	return t
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any any) time.Duration {
	d, _ := defaultConverter.Duration(any)
	return d
}

// GTime converts `any` to *time.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using jtime.NewFromTimeStamp if `any` is numeric,
// or using jtime.StrToTime if `any` is string.
func GTime(any any, format ...string) *jtime.Time {
	t, _ := defaultConverter.GTime(any, format...)
	return t
}
