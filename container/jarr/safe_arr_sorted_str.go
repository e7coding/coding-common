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
	"github.com/e7coding/coding-common/text/jstr"
	"sort"
	"strings"
	"sync"
)

// SafeSortedStrArr 有序字符串数组，支持去重和自定义比较
type SafeSortedStrArr struct {
	mu     sync.RWMutex
	data   []string
	unique bool
	cmp    func(a, b string) int
}

// NewSafeSortedStrArr 创建空数组，可指定比较函数
func NewSafeSortedStrArr(cmp ...func(a, b string) int) *SafeSortedStrArr {
	f := CompareStrings
	if len(cmp) > 0 && cmp[0] != nil {
		f = cmp[0]
	}
	return &SafeSortedStrArr{data: []string{}, cmp: f}
}

// NewSafeSortedStrArrWithCap 指定容量创建数组
func NewSafeSortedStrArrWithCap(cap int, cmp ...func(a, b string) int) *SafeSortedStrArr {
	a := NewSafeSortedStrArr(cmp...)
	a.data = make([]string, 0, cap)
	return a
}

// NewSafeSortedStrFrom 根据切片创建并排序
func NewSafeSortedStrFrom(src []string, cmp ...func(a, b string) int) *SafeSortedStrArr {
	a := NewSafeSortedStrArrWithCap(len(src), cmp...)
	a.data = append(a.data, src...)
	return a.Sort()
}

// Append 添加元素并排序
func (a *SafeSortedStrArr) Append(vals ...string) *SafeSortedStrArr {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data = append(a.data, vals...)
	return a.Sort()
}

// Add Append
func (a *SafeSortedStrArr) Add(vals ...string) *SafeSortedStrArr {
	return a.Append(vals...)
}

// Sort 排序并去重（若已设置去重）
func (a *SafeSortedStrArr) Sort() *SafeSortedStrArr {
	a.mu.Lock()
	defer a.mu.Unlock()
	sort.Slice(a.data, func(i, j int) bool {
		return a.cmp(a.data[i], a.data[j]) < 0
	})
	if a.unique {
		a.uniqFilter()
	}
	return a
}

// SetUniq 设置去重
func (a *SafeSortedStrArr) SetUniq(u bool) *SafeSortedStrArr {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.unique = u
	if u {
		a.uniqFilter()
	}
	return a
}

// uniqFilter 执行去重
func (a *SafeSortedStrArr) uniqFilter() {
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

// At 获取元素，越界返回空串
func (a *SafeSortedStrArr) At(i int) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if i < 0 || i >= len(a.data) {
		return ""
	}
	return a.data[i]
}

// Len 数量
func (a *SafeSortedStrArr) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.data)
}

// Slice 底层切片
func (a *SafeSortedStrArr) Slice() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.data
}

// Sum 元素转整数后求和
func (a *SafeSortedStrArr) Sum() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	s := 0
	for _, v := range a.data {
		s += jconv.Int(v)
	}
	return s
}

// Contains 是否包含
func (a *SafeSortedStrArr) Contains(v string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Search(v) != -1
}

// ContainsI 忽略大小写包含检查
func (a *SafeSortedStrArr) ContainsI(v string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, s := range a.data {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}

// Search 二分查找，找不到返回-1
func (a *SafeSortedStrArr) Search(v string) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	i, r := a.binarySearch(v)
	if r == 0 {
		return i
	}
	return -1
}

// binarySearch 二分查找，返回位置及比较结果
func (a *SafeSortedStrArr) binarySearch(v string) (idx, res int) {
	lo, hi := 0, len(a.data)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		res = a.cmp(v, a.data[mid])
		if res < 0 {
			hi = mid - 1
		} else if res > 0 {
			lo = mid + 1
		} else {
			return mid, 0
		}
	}
	return lo, res
}

// PopLeft 弹出首元素
func (a *SafeSortedStrArr) PopLeft() (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.data) == 0 {
		return "", false
	}
	v := a.data[0]
	a.data = a.data[1:]
	return v, true
}

// PopRight 弹出尾元素
func (a *SafeSortedStrArr) PopRight() (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	L := len(a.data)
	if L == 0 {
		return "", false
	}
	v := a.data[L-1]
	a.data = a.data[:L-1]
	return v, true
}

// PopRand 随机弹出
func (a *SafeSortedStrArr) PopRand() (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.data) == 0 {
		return "", false
	}
	pos := jrand.Intn(len(a.data))
	v := a.data[pos]
	a.data = append(a.data[:pos], a.data[pos+1:]...)
	return v, true
}

// Join 拼接字符串
func (a *SafeSortedStrArr) Join(sep string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	buf := &bytes.Buffer{}
	for i, v := range a.data {
		if jstr.IsNumeric(v) {
			buf.WriteString(v)
		} else {
			buf.WriteString(`"` + jstr.QuoteMeta(v, `"\`) + `"`)
		}
		if i < len(a.data)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}

// MarshalJSON 自定义序列化
func (a *SafeSortedStrArr) MarshalJSON() ([]byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return json.Marshal(a.data)
}
