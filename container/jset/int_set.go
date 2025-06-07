package jset

import (
	"bytes"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// IntSet 整型集合，内部使用 map[int]struct{} 存储元素
type IntSet struct {
	data map[int]struct{}
}

// NewIntSet 创建并返回一个空的 IntSet
func NewIntSet() *IntSet {
	return &IntSet{data: make(map[int]struct{})}
}

// NewIntSetFrom 从给定的切片初始化并返回一个 IntSet
func NewIntSetFrom(items []int) *IntSet {
	m := make(map[int]struct{}, len(items))
	for _, v := range items {
		m[v] = struct{}{}
	}
	return &IntSet{data: m}
}

// Iterator 遍历集合中的元素，f 返回 false 时停止
func (s *IntSet) Iterator(f func(v int) bool) {
	for k := range s.data {
		if !f(k) {
			break
		}
	}
}

// Add 向集合中添加一个或多个元素
func (s *IntSet) Add(items ...int) {
	if s.data == nil {
		s.data = make(map[int]struct{})
	}
	for _, v := range items {
		s.data[v] = struct{}{}
	}
}

// Contains 判断集合是否包含指定元素
func (s *IntSet) Contains(item int) bool {
	_, ok := s.data[item]
	return ok
}

// Remove 从集合中删除指定元素
func (s *IntSet) Remove(item int) {
	delete(s.data, item)
}

// Size 返回集合中元素的数量
func (s *IntSet) Size() int {
	return len(s.data)
}

// Clear 清空集合
func (s *IntSet) Clear() {
	s.data = make(map[int]struct{})
}

// Slice 返回集合元素的切片
func (s *IntSet) Slice() []int {
	ret := make([]int, 0, len(s.data))
	for k := range s.data {
		ret = append(ret, k)
	}
	return ret
}

// Join 将集合元素转换为字符串并以 glue 分隔
func (s *IntSet) Join(glue string) string {
	if len(s.data) == 0 {
		return ""
	}
	var buf = bytes.NewBuffer(nil)
	i, n := 0, len(s.data)
	for k := range s.data {
		buf.WriteString(jconv.String(k))
		if i < n-1 {
			buf.WriteString(glue)
		}
		i++
	}
	return buf.String()
}

// String 返回集合的 JSON 样式表示，如 [1,2,3]
func (s *IntSet) String() string {
	return "[" + s.Join(",") + "]"
}

// Equal 判断两个集合是否相等
func (s *IntSet) Equal(o *IntSet) bool {
	if s == o {
		return true
	}
	if len(s.data) != len(o.data) {
		return false
	}
	for k := range s.data {
		if _, ok := o.data[k]; !ok {
			return false
		}
	}
	return true
}

// IsSubsetOf 判断当前集合是否是 o 的子集
func (s *IntSet) IsSubsetOf(o *IntSet) bool {
	if s == o {
		return true
	}
	for k := range s.data {
		if _, ok := o.data[k]; !ok {
			return false
		}
	}
	return true
}

// Union 返回当前集合与多个其他集合的并集
func (s *IntSet) Union(others ...*IntSet) *IntSet {
	res := NewIntSet()
	for k := range s.data {
		res.data[k] = struct{}{}
	}
	for _, o := range others {
		for k := range o.data {
			res.data[k] = struct{}{}
		}
	}
	return res
}

// Diff 返回当前集合与多个其他集合的差集（在 s 中但不在 o 中）
func (s *IntSet) Diff(others ...*IntSet) *IntSet {
	res := NewIntSet()
next:
	for k := range s.data {
		for _, o := range others {
			if _, ok := o.data[k]; ok {
				continue next
			}
		}
		res.data[k] = struct{}{}
	}
	return res
}

// Intersect 返回当前集合与多个其他集合的交集
func (s *IntSet) Intersect(others ...*IntSet) *IntSet {
	res := NewIntSet()
	for k := range s.data {
		ok := true
		for _, o := range others {
			if _, found := o.data[k]; !found {
				ok = false
				break
			}
		}
		if ok {
			res.data[k] = struct{}{}
		}
	}
	return res
}

// Complement 返回 full 集合中不在当前集合 s 中的元素
func (s *IntSet) Complement(full *IntSet) *IntSet {
	res := NewIntSet()
	for k := range full.data {
		if _, ok := s.data[k]; !ok {
			res.data[k] = struct{}{}
		}
	}
	return res
}

// Merge 将多个集合的元素添加到当前集合并返回自身
func (s *IntSet) Merge(others ...*IntSet) *IntSet {
	for _, o := range others {
		for k := range o.data {
			s.data[k] = struct{}{}
		}
	}
	return s
}

// Sum 对集合中的元素求和并返回结果
func (s *IntSet) Sum() int {
	sum := 0
	for k := range s.data {
		sum += k
	}
	return sum
}

// Pop 随机弹出并删除集合中的一个元素，若集合为空返回 0
func (s *IntSet) Pop() int {
	for k := range s.data {
		delete(s.data, k)
		return k
	}
	return 0
}

// Pops 随机弹出并删除 size 个元素，size<0 或 超过长度时弹出所有
func (s *IntSet) Pops(size int) []int {
	n := len(s.data)
	if size < 0 || size > n {
		size = n
	}
	res := make([]int, 0, size)
	i := 0
	for k := range s.data {
		if i >= size {
			break
		}
		delete(s.data, k)
		res = append(res, k)
		i++
	}
	return res
}

// Walk 对集合中的每个元素应用 f 并返回新集合
func (s *IntSet) Walk(f func(int) int) *IntSet {
	res := NewIntSet()
	for k := range s.data {
		res.data[f(k)] = struct{}{}
	}
	return res
}

// MarshalJSON 将集合编码为 JSON 数组
func (s *IntSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}

// UnmarshalJSON 从 JSON 数组解码到集合
func (s *IntSet) UnmarshalJSON(b []byte) error {
	var arr []int
	if err := json.UnmarshalUseNumber(b, &arr); err != nil {
		return err
	}
	s.data = make(map[int]struct{}, len(arr))
	for _, v := range arr {
		s.data[v] = struct{}{}
	}
	return nil
}

// UnmarshalValue 从任意类型（切片或 JSON）解码到集合
func (s *IntSet) UnmarshalValue(value interface{}) error {
	arr := jconv.SliceInt(value)
	s.data = make(map[int]struct{}, len(arr))
	for _, v := range arr {
		s.data[v] = struct{}{}
	}
	return nil
}

// DeepCopy 返回集合的深拷贝
func (s *IntSet) DeepCopy() interface{} {
	return NewIntSetFrom(s.Slice())
}
