// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jconv

import (
	"github.com/e7coding/coding-common/jutil/jconv/internal/converter"
)

// SliceMap is alias of Maps.
func SliceMap(any any, option ...MapOption) []map[string]any {
	return Maps(any, option...)
}

// Maps converts `value` to []map[string]any.
// Note that it automatically checks and converts json string to []map if `value` is string/[]byte.
func Maps(value any, option ...MapOption) []map[string]any {
	mapOption := MapOption{
		ContinueOnError: true,
	}
	if len(option) > 0 {
		mapOption = option[0]
	}
	result, _ := defaultConverter.SliceMap(value, SliceMapOption{
		MapOption: mapOption,
		SliceOption: converter.SliceOption{
			ContinueOnError: true,
		},
	})
	return result
}
