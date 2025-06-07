// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jarr

import "strings"

// ComparatorInt 比较整数 a, b 大小
type ComparatorInt func(a, b int) int

// ComparatorString 比较字符串 a, b 大小
type ComparatorString func(a, b string) int

// compareInt 默认整数比较，a<b返回-1，a>b返回1，否则返回0
func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// compareString 默认字符串比较，基于 Unicode 编码
func compareString(a, b string) int {
	return strings.Compare(a, b)
}

// quickSortInt 快速排序整型切片
func quickSortInt(vals []int, cmp ComparatorInt) {
	if len(vals) <= 1 {
		return
	}
	pivot := vals[len(vals)/2]
	i, j := 0, len(vals)-1
	for i <= j {
		for cmp(vals[i], pivot) < 0 {
			i++
		}
		for cmp(vals[j], pivot) > 0 {
			j--
		}
		if i <= j {
			vals[i], vals[j] = vals[j], vals[i]
			i++
			j--
		}
	}
	if j > 0 {
		quickSortInt(vals[:j+1], cmp)
	}
	if i < len(vals)-1 {
		quickSortInt(vals[i:], cmp)
	}
}

// quickSortString 快速排序字符串切片
func quickSortString(vals []string, cmp ComparatorString) {
	if len(vals) <= 1 {
		return
	}
	pivot := vals[len(vals)/2]
	i, j := 0, len(vals)-1
	for i <= j {
		for cmp(vals[i], pivot) < 0 {
			i++
		}
		for cmp(vals[j], pivot) > 0 {
			j--
		}
		if i <= j {
			vals[i], vals[j] = vals[j], vals[i]
			i++
			j--
		}
	}
	if j > 0 {
		quickSortString(vals[:j+1], cmp)
	}
	if i < len(vals)-1 {
		quickSortString(vals[i:], cmp)
	}
}
