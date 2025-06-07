// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jarr

import (
	"bytes"
	"encoding/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/jutil/jrand"
	"sort"
)

// SortedIntArray 有序整数数组
// 支持去重和自定义比较
type SortedIntArray struct {
	data   []int
	unique bool
	cmp    func(a, b int) int
}

// CompareInts 默认整数比较
func CompareInts(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// NewSortedIntArray 创建空数组，可选自定义比较
func NewSortedIntArray(cmp ...func(a, b int) int) *SortedIntArray {
	c := CompareInts
	if len(cmp) > 0 && cmp[0] != nil {
		c = cmp[0]
	}
	return &SortedIntArray{data: []int{}, cmp: c}
}

// NewWithCap 创建指定容量数组，可选自定义比较
func NewWithCap(cap int, cmp ...func(a, b int) int) *SortedIntArray {
	a := NewSortedIntArray(cmp...)
	a.data = make([]int, 0, cap)
	return a
}

// NewFrom 根据已有切片创建并排序，可选自定义比较
func NewFrom(src []int, cmp ...func(a, b int) int) *SortedIntArray {
	a := NewWithCap(len(src), cmp...)
	a.data = append(a.data, src...)
	return a.Sort()
}

// Append 添加元素并排序
func (a *SortedIntArray) Append(vals ...int) *SortedIntArray {
	a.data = append(a.data, vals...)
	return a.Sort()
}

// Add 为 Append 别名
func (a *SortedIntArray) Add(vals ...int) *SortedIntArray {
	return a.Append(vals...)
}

// Sort 排序并去重（若设置去重）
func (a *SortedIntArray) Sort() *SortedIntArray {
	sort.Slice(a.data, func(i, j int) bool {
		return a.cmp(a.data[i], a.data[j]) < 0
	})
	if a.unique {
		a.uniqueFilter()
	}
	return a
}

// SetUnique 设置去重
func (a *SortedIntArray) SetUnique(u bool) *SortedIntArray {
	a.unique = u
	if u {
		a.uniqueFilter()
	}
	return a
}

// uniqueFilter 执行去重
func (a *SortedIntArray) uniqueFilter() {
	if len(a.data) < 2 {
		return
	}
	j := 0
	for i := 1; i < len(a.data); i++ {
		if a.cmp(a.data[j], a.data[i]) != 0 {
			j++
			a.data[j] = a.data[i]
		}
	}
	a.data = a.data[:j+1]
}

// At 获取索引元素，越界返回0
func (a *SortedIntArray) At(idx int) int {
	if idx < 0 || idx >= len(a.data) {
		return 0
	}
	return a.data[idx]
}

// Len 长度
func (a *SortedIntArray) Len() int {
	return len(a.data)
}

// Slice 底层切片
func (a *SortedIntArray) Slice() []int {
	return a.data
}

// Sum 求和
func (a *SortedIntArray) Sum() int {
	s := 0
	for _, v := range a.data {
		s += v
	}
	return s
}

// Contains 是否包含
func (a *SortedIntArray) Contains(v int) bool {
	return a.Search(v) != -1
}

// Search 查找索引，找不到返回-1
func (a *SortedIntArray) Search(v int) int {
	i, r := a.binarySearch(v)
	if r == 0 {
		return i
	}
	return -1
}

// binarySearch 二分查找，返回位置和比较结果
func (a *SortedIntArray) binarySearch(v int) (idx, res int) {
	low, high := 0, len(a.data)-1
	res = -2
	for low <= high {
		mid := (low + high) / 2
		res = a.cmp(v, a.data[mid])
		if res < 0 {
			high = mid - 1
		} else if res > 0 {
			low = mid + 1
		} else {
			return mid, res
		}
	}
	return low, res
}

// PopLeft 弹出第一个元素
func (a *SortedIntArray) PopLeft() (int, bool) {
	if len(a.data) == 0 {
		return 0, false
	}
	v := a.data[0]
	a.data = a.data[1:]
	return v, true
}

// PopRight 弹出最后一个元素
func (a *SortedIntArray) PopRight() (int, bool) {
	if len(a.data) == 0 {
		return 0, false
	}
	idx := len(a.data) - 1
	v := a.data[idx]
	a.data = a.data[:idx]
	return v, true
}

// PopRand 随机弹出一个元素
func (a *SortedIntArray) PopRand() (int, bool) {
	if len(a.data) == 0 {
		return 0, false
	}
	pos := jrand.Intn(len(a.data))
	v := a.data[pos]
	a.data = append(a.data[:pos], a.data[pos+1:]...)
	return v, true
}

// Join 拼接为字符串
func (a *SortedIntArray) Join(glue string) string {
	buf := &bytes.Buffer{}
	for i, v := range a.data {
		buf.WriteString(jconv.String(v))
		if i < len(a.data)-1 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// MarshalJSON 自定义 JSON 序列化
func (a *SortedIntArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.data)
}
