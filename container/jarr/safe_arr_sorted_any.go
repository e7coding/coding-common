// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jarr

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"sync"
)

// SafeSortedArr 有序数组，支持并发安全和去重
type SafeSortedArr struct {
	mu      sync.RWMutex
	items   []interface{}
	unique  bool
	compare ComparatorSortedArr
}

// NewSafeSortedArr 创建空有序数组，需传入比较函数
func NewSafeSortedArr(cmp ComparatorSortedArr) *SafeSortedArr {
	return &SafeSortedArr{items: []interface{}{}, compare: cmp}
}

// NewSafeSortedArrFrom 从切片创建并排序
func NewSafeSortedArrFrom(slice []interface{}, cmp ComparatorSortedArr) *SafeSortedArr {
	sa := NewSafeSortedArr(cmp)
	sa.items = append(sa.items, slice...)
	sa.Sort()
	return sa
}

// Len 返回长度
func (sa *SafeSortedArr) Len() int {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	n := len(sa.items)

	return n
}

// At 返回索引处元素，越界返回 error
func (sa *SafeSortedArr) At(i int) (interface{}, error) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	if i < 0 || i >= len(sa.items) {

		return nil, errors.New("index out of range")
	}
	v := sa.items[i]

	return v, nil
}

// Slice 返回数据副本
func (sa *SafeSortedArr) Slice() []interface{} {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	copyItems := make([]interface{}, len(sa.items))
	copy(copyItems, sa.items)

	return copyItems
}

// Sort 排序并去重（若已开启去重）
func (sa *SafeSortedArr) Sort() {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sort.Slice(sa.items, func(i, j int) bool {
		return sa.compare(sa.items[i], sa.items[j]) < 0
	})
	if sa.unique {
		sa.uniqAllLocked()
	}

}

// Add 插入元素并保持排序
func (sa *SafeSortedArr) Add(vals ...interface{}) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	for _, v := range vals {
		idx, eq := sa.binarySearchLocked(v)
		if eq == 0 && sa.unique {
			continue
		}
		if eq > 0 {
			idx++
		}
		sa.items = append(sa.items, nil)
		copy(sa.items[idx+1:], sa.items[idx:])
		sa.items[idx] = v
	}

}

// Remove 删除索引处元素
func (sa *SafeSortedArr) Remove(i int) (interface{}, error) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	if i < 0 || i >= len(sa.items) {

		return nil, errors.New("index out of range")
	}
	v := sa.items[i]
	sa.items = append(sa.items[:i], sa.items[i+1:]...)

	return v, nil
}

// Search 二分查找，找不到返回 -1
func (sa *SafeSortedArr) Search(v interface{}) int {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	i, eq := sa.binarySearchLocked(v)

	if eq == 0 {
		return i
	}
	return -1
}

// Contains 是否包含
func (sa *SafeSortedArr) Contains(v interface{}) bool {
	return sa.Search(v) != -1
}

// SetUnique 设置去重
func (sa *SafeSortedArr) SetUnique(u bool) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.unique = u
	if u {
		sa.uniqAllLocked()
	}

}

// PopFront 弹出首元素
func (sa *SafeSortedArr) PopFront() (interface{}, error) {
	return sa.Remove(0)
}

// PopBack 弹出尾元素
func (sa *SafeSortedArr) PopBack() (interface{}, error) {
	return sa.Remove(len(sa.items) - 1)
}

// PopRandom 弹出随机元素
func (sa *SafeSortedArr) PopRandom() (interface{}, error) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	n := len(sa.items)
	if n == 0 {

		return nil, errors.New("empty array")
	}
	i := rand.Intn(n)
	v := sa.items[i]
	sa.items = append(sa.items[:i], sa.items[i+1:]...)

	return v, nil
}

// Clear 清空
func (sa *SafeSortedArr) Clear() {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.items = sa.items[:0]

}

// Clone 克隆副本
func (sa *SafeSortedArr) Clone() *SafeSortedArr {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	copyItems := make([]interface{}, len(sa.items))
	copy(copyItems, sa.items)

	return NewSafeSortedArrFrom(copyItems, sa.compare)
}

// MarshalJSON 自定义序列化
func (sa *SafeSortedArr) MarshalJSON() ([]byte, error) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	data := sa.items

	return json.Marshal(data)
}

// UnmarshalJSON 自定义反序列化
func (sa *SafeSortedArr) UnmarshalJSON(b []byte) error {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	if sa.compare == nil {

		return errors.New("comparator not set")
	}
	var a []interface{}
	if err := json.Unmarshal(b, &a); err != nil {

		return err
	}
	sa.items = a
	sort.Slice(sa.items, func(i, j int) bool {
		return sa.compare(sa.items[i], sa.items[j]) < 0
	})
	if sa.unique {
		sa.uniqAllLocked()
	}

	return nil
}

// uniqAllLocked 去重，需持写锁
func (sa *SafeSortedArr) uniqAllLocked() {
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

// binarySearchLocked 二分查找，需持读锁
func (sa *SafeSortedArr) binarySearchLocked(v interface{}) (idx, cmp int) {
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
