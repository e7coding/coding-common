// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jarr

import (
	"encoding/json"
	"errors"
	"github.com/e7coding/coding-common/jutil/jrand"
	"sort"
)

// ComparatorSortedArr 定义比较函数
// 返回值 <0: a<b, 0: a==b, >0: a>b
type ComparatorSortedArr func(a, b interface{}) int

// SortedArray 有序数组
type SortedArray struct {
	items   []interface{}
	unique  bool
	compare ComparatorSortedArr
}

// NewSortedArray 创建空有序数组
func NewSortedArray(cmp ComparatorSortedArr) *SortedArray {
	return &SortedArray{items: []interface{}{}, compare: cmp}
}

// NewSortedArrayFrom 从切片创建有序数组
func NewSortedArrayFrom(slice []interface{}, cmp ComparatorSortedArr) *SortedArray {
	sa := NewSortedArray(cmp)
	sa.items = append(sa.items, slice...)
	sa.Sort()
	return sa
}

// Len 返回长度
func (sa *SortedArray) Len() int {
	return len(sa.items)
}

// At 返回索引处的元素，越界返回 error
func (sa *SortedArray) At(i int) (interface{}, error) {
	if i < 0 || i >= len(sa.items) {
		return nil, errors.New("index out of range")
	}
	return sa.items[i], nil
}

// Slice 返回底层切片
func (sa *SortedArray) Slice() []interface{} {
	return sa.items
}

// Sort 排序
func (sa *SortedArray) Sort() {
	sort.Slice(sa.items, func(i, j int) bool {
		return sa.compare(sa.items[i], sa.items[j]) < 0
	})
}

// Add 插入元素
func (sa *SortedArray) Add(vals ...interface{}) {
	for _, v := range vals {
		idx, eq := sa.binarySearch(v)
		if eq == 0 && sa.unique {
			continue
		}
		// 插入位置
		if eq > 0 {
			idx++
		}
		sa.items = append(sa.items, nil)
		copy(sa.items[idx+1:], sa.items[idx:])
		sa.items[idx] = v
	}
}

// Remove 删除索引处元素，返回元素或 error
func (sa *SortedArray) Remove(i int) (interface{}, error) {
	if i < 0 || i >= len(sa.items) {
		return nil, errors.New("index out of range")
	}
	v := sa.items[i]
	sa.items = append(sa.items[:i], sa.items[i+1:]...)
	return v, nil
}

// Search 二分查找，返回索引或 -1
func (sa *SortedArray) Search(v interface{}) int {
	idx, eq := sa.binarySearch(v)
	if eq == 0 {
		return idx
	}
	return -1
}

// Contains 判断是否存在
func (sa *SortedArray) Contains(v interface{}) bool {
	return sa.Search(v) != -1
}

// SetUniq 设置去重标志
func (sa *SortedArray) SetUniq(u bool) {
	sa.unique = u
	if u {
		sa.uniqAll()
	}
}

// uniqAll 去重
func (sa *SortedArray) uniqAll() {
	if len(sa.items) < 2 {
		return
	}
	n := sa.items[:1]
	for _, v := range sa.items[1:] {
		if sa.compare(n[len(n)-1], v) != 0 {
			n = append(n, v)
		}
	}
	sa.items = n
}

// PopFront 弹出第一个元素
func (sa *SortedArray) PopFront() (interface{}, error) {
	return sa.Remove(0)
}

// PopBack 弹出最后一个元素
func (sa *SortedArray) PopBack() (interface{}, error) {
	return sa.Remove(len(sa.items) - 1)
}

// PopRandom 弹出随机元素
func (sa *SortedArray) PopRandom() (interface{}, error) {
	if sa.Len() == 0 {
		return nil, errors.New("empty array")
	}
	i := jrand.Intn(sa.Len())
	return sa.Remove(i)
}

// ForEach 遍历 , 回调返回 false 停止
func (sa *SortedArray) ForEach(f func(index int, v interface{}) bool) {
	for i, v := range sa.items {
		if !f(i, v) {
			break
		}
	}
}

// ForEachReverse 逆序遍历
func (sa *SortedArray) ForEachReverse(f func(index int, v interface{}) bool) {
	for i := len(sa.items) - 1; i >= 0; i-- {
		if !f(i, sa.items[i]) {
			break
		}
	}
}

// Sum 返回整数总和
func (sa *SortedArray) Sum() int {
	s, tot := sa.items, 0
	for _, v := range s {
		if n, ok := v.(int); ok {
			tot += n
		}
	}
	return tot
}

// Range 返回[start,end)子切片
func (sa *SortedArray) Range(start, end int) []interface{} {
	if start < 0 || end > len(sa.items) || start >= end {
		return nil
	}
	return sa.items[start:end]
}

// Clear 清空
func (sa *SortedArray) Clear() {
	sa.items = sa.items[:0]
}

// Clone 克隆
func (sa *SortedArray) Clone() *SortedArray {
	n := make([]interface{}, len(sa.items))
	copy(n, sa.items)
	return NewSortedArrayFrom(n, sa.compare)
}

// MarshalJSON JSON 序列化
func (sa *SortedArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(sa.items)
}

// UnmarshalJSON JSON 反序列化
func (sa *SortedArray) UnmarshalJSON(b []byte) error {
	if sa.compare == nil {
		return errors.New("comparator not set")
	}
	a := []interface{}{}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	sa.items = a
	sa.Sort()
	return nil
}

// binarySearch 私有二分查找
// 返回位置及比较结果
func (sa *SortedArray) binarySearch(v interface{}) (index, cmp int) {
	low, high := 0, len(sa.items)-1
	for low <= high {
		mid := (low + high) / 2
		cmp = sa.compare(v, sa.items[mid])
		switch {
		case cmp < 0:
			high = mid - 1
		case cmp > 0:
			low = mid + 1
		default:
			return mid, 0
		}
	}
	return low, cmp
}
