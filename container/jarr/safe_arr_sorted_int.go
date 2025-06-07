// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jarr

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/e7coding/coding-common/jutil/jconv"
	"math/rand"
	"sort"
	"sync"
)

// SafeSortedIntArr 有序整数数组，支持并发安全和去重
type SafeSortedIntArr struct {
	mu     sync.RWMutex
	data   []int
	unique bool
	cmp    func(a, b int) int
}

// NewSafeSortedIntArr 创建空数组，可选自定义比较
func NewSafeSortedIntArr(cmp ...func(a, b int) int) *SafeSortedIntArr {
	c := CompareInts
	if len(cmp) > 0 && cmp[0] != nil {
		c = cmp[0]
	}
	return &SafeSortedIntArr{data: []int{}, cmp: c}
}

// NewSafeSortedIntArrWithCap 创建指定容量数组
func NewSafeSortedIntArrWithCap(capacity int, cmp ...func(a, b int) int) *SafeSortedIntArr {
	a := NewSafeSortedIntArr(cmp...)
	a.data = make([]int, 0, capacity)
	return a
}

// NewSafeSortedIntArrFrom 根据切片创建并排序
func NewSafeSortedIntArrFrom(src []int, cmp ...func(a, b int) int) *SafeSortedIntArr {
	a := NewSafeSortedIntArrWithCap(len(src), cmp...)
	a.data = append(a.data, src...)
	a.Sort()
	return a
}

// Len 长度
func (a *SafeSortedIntArr) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	n := len(a.data)

	return n
}

// At 获取索引元素，越界报错
func (a *SafeSortedIntArr) At(idx int) (int, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if idx < 0 || idx >= len(a.data) {

		return 0, errors.New("index out of range")
	}
	v := a.data[idx]

	return v, nil
}

// Slice 返回数据副本
func (a *SafeSortedIntArr) Slice() []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]int, len(a.data))
	copy(out, a.data)

	return out
}

// Sort 排序并去重
func (a *SafeSortedIntArr) Sort() *SafeSortedIntArr {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.data, func(i, j int) bool {
		return a.cmp(a.data[i], a.data[j]) < 0
	})
	if a.unique {
		a.uniqueFilter()
	}

	return a
}

// SetUnique 设置去重
func (a *SafeSortedIntArr) SetUnique(u bool) *SafeSortedIntArr {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.unique = u
	if u {
		a.uniqueFilter()
	}

	return a
}

// uniqueFilter 去重，需持写锁
func (a *SafeSortedIntArr) uniqueFilter() {
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

// Append 添加元素并排序
func (a *SafeSortedIntArr) Append(vals ...int) *SafeSortedIntArr {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data = append(a.data, vals...)

	return a.Sort()
}

// Add Append 别名
func (a *SafeSortedIntArr) Add(vals ...int) *SafeSortedIntArr {
	return a.Append(vals...)
}

// Search 二分查找，找不到返回-1
func (a *SafeSortedIntArr) Search(v int) int {
	a.mu.RLock()
	a.mu.RUnlock()
	idx, res := a.binarySearch(v)

	if res == 0 {
		return idx
	}
	return -1
}

// Contains 是否包含
func (a *SafeSortedIntArr) Contains(v int) bool {
	return a.Search(v) != -1
}

// binarySearch 私有二分查找
func (a *SafeSortedIntArr) binarySearch(v int) (idx, res int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	low, high := 0, len(a.data)-1
	for low <= high {
		mid := (low + high) / 2
		res = a.cmp(v, a.data[mid])
		if res < 0 {
			high = mid - 1
		} else if res > 0 {
			low = mid + 1
		} else {
			return mid, 0
		}
	}
	return low, res
}

// PopLeft 弹出首元素
func (a *SafeSortedIntArr) PopLeft() (int, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.data) == 0 {

		return 0, false
	}
	v := a.data[0]
	a.data = a.data[1:]

	return v, true
}

// PopRight 弹出尾元素
func (a *SafeSortedIntArr) PopRight() (int, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	L := len(a.data)
	if L == 0 {

		return 0, false
	}
	v := a.data[L-1]
	a.data = a.data[:L-1]

	return v, true
}

// PopRand 随机弹出
func (a *SafeSortedIntArr) PopRand() (int, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	L := len(a.data)
	if L == 0 {

		return 0, false
	}
	pos := rand.Intn(L)
	v := a.data[pos]
	a.data = append(a.data[:pos], a.data[pos+1:]...)

	return v, true
}

// Join 拼接字符串
func (a *SafeSortedIntArr) Join(glue string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	buf := &bytes.Buffer{}
	for i, v := range a.data {
		buf.WriteString(jconv.String(v))
		if i < len(a.data)-1 {
			buf.WriteString(glue)
		}
	}

	return buf.String()
}

// MarshalJSON 自定义序列化
func (a *SafeSortedIntArr) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	a.mu.RUnlock()
	data := a.data

	return json.Marshal(data)
}
