package jset

import (
	"bytes"
	"strings"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// StrSet 字符串集合
type StrSet struct {
	data map[string]struct{}
}

// NewStrSet 创建并返回一个空的字符串集合
func NewStrSet() *StrSet {
	return &StrSet{data: make(map[string]struct{})}
}

// NewStrSetFrom 根据给定的切片创建并返回一个字符串集合
func NewStrSetFrom(items []string) *StrSet {
	m := make(map[string]struct{}, len(items))
	for _, v := range items {
		m[v] = struct{}{}
	}
	return &StrSet{data: m}
}

// Foreach 迭代集合中的每个元素，f 返回 false 时停止
func (s *StrSet) Foreach(f func(v string) bool) {
	for k := range s.data {
		if !f(k) {
			break
		}
	}
}

// Add 向集合中添加一个或多个元素
func (s *StrSet) Add(items ...string) {
	if s.data == nil {
		s.data = make(map[string]struct{})
	}
	for _, v := range items {
		s.data[v] = struct{}{}
	}
}

// AddIfNotExist 如果元素不存在则添加并返回 true，否则返回 false
func (s *StrSet) AddIfNotExist(item string) bool {
	if !s.Contains(item) {
		if s.data == nil {
			s.data = make(map[string]struct{})
		}
		s.data[item] = struct{}{}
		return true
	}
	return false
}

func (s *StrSet) AddIfNotExistFunc(item string, f func() bool) bool {
	if !s.Contains(item) {
		if f() {
			if s.data == nil {
				s.data = make(map[string]struct{})
			}
			if _, ok := s.data[item]; !ok {
				s.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// Contains 判断集合是否包含指定元素
func (s *StrSet) Contains(item string) bool {
	_, ok := s.data[item]
	return ok
}

// ContainsI 不区分大小写地判断元素是否存在
func (s *StrSet) ContainsI(item string) bool {
	for k := range s.data {
		if strings.EqualFold(k, item) {
			return true
		}
	}
	return false
}

// Remove 从集合中删除指定元素
func (s *StrSet) Remove(item string) {
	delete(s.data, item)
}

// Size 返回集合中元素的数量
func (s *StrSet) Size() int {
	return len(s.data)
}

// Clear 清空集合
func (s *StrSet) Clear() {
	s.data = make(map[string]struct{})
}

// Slice 返回集合元素的切片
func (s *StrSet) Slice() []string {
	ret := make([]string, 0, len(s.data))
	for k := range s.data {
		ret = append(ret, k)
	}
	return ret
}

// Join 将集合元素以 glue 连接成字符串
func (s *StrSet) Join(glue string) string {
	if len(s.data) == 0 {
		return ""
	}
	var buf = bytes.NewBuffer(nil)
	i, n := 0, len(s.data)
	for k := range s.data {
		buf.WriteString(k)
		if i < n-1 {
			buf.WriteString(glue)
		}
		i++
	}
	return buf.String()
}

// String 返回集合的 JSON 样式字符串表示
func (s *StrSet) String() string {
	return `[` + s.Join(`","`) + `]`
}

// Equal 判断两个集合是否相等
func (s *StrSet) Equal(o *StrSet) bool {
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

// IsSubsetOf 判断当前集合是否为 o 的子集
func (s *StrSet) IsSubsetOf(o *StrSet) bool {
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

// Union 返回当前集合与多个集合的并集
func (s *StrSet) Union(others ...*StrSet) *StrSet {
	res := NewStrSet()
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

// Diff 返回当前集合与多个集合的差集（在 s 中但不在 others 中）
func (s *StrSet) Diff(others ...*StrSet) *StrSet {
	res := NewStrSet()
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

// Intersect 返回当前集合与多个集合的交集
func (s *StrSet) Intersect(others ...*StrSet) *StrSet {
	res := NewStrSet()
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

// Complement 返回 full 集合中不在当前集合中的元素
func (s *StrSet) Complement(full *StrSet) *StrSet {
	res := NewStrSet()
	for k := range full.data {
		if _, ok := s.data[k]; !ok {
			res.data[k] = struct{}{}
		}
	}
	return res
}

// Merge 将多个集合的元素合并到当前集合并返回自身
func (s *StrSet) Merge(others ...*StrSet) *StrSet {
	for _, o := range others {
		for k := range o.data {
			s.data[k] = struct{}{}
		}
	}
	return s
}

// Sum 对集合中可转换为 int 的元素求和
func (s *StrSet) Sum() int {
	sum := 0
	for k := range s.data {
		sum += jconv.Int(k)
	}
	return sum
}

// Pop 随机弹出并删除一个元素，若空则返回空串
func (s *StrSet) Pop() string {
	for k := range s.data {
		delete(s.data, k)
		return k
	}
	return ""
}

// Pops 随机弹出并删除 size 个元素，size<0 或 超过长度时弹出所有
func (s *StrSet) Pops(size int) []string {
	n := len(s.data)
	if size < 0 || size > n {
		size = n
	}
	res := make([]string, 0, size)
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

// Walk 对集合中每个元素应用 f 并返回新集合
func (s *StrSet) Walk(f func(string) string) *StrSet {
	res := NewStrSet()
	for k := range s.data {
		res.data[f(k)] = struct{}{}
	}
	return res
}

// MarshalJSON 将集合编码为 JSON 数组
func (s *StrSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}

// UnmarshalJSON 从 JSON 数组解码到集合
func (s *StrSet) UnmarshalJSON(b []byte) error {
	var arr []string
	if err := json.UnmarshalUseNumber(b, &arr); err != nil {
		return err
	}
	s.data = make(map[string]struct{}, len(arr))
	for _, v := range arr {
		s.data[v] = struct{}{}
	}
	return nil
}

// UnmarshalValue 从任意类型（切片或 JSON）解码到集合
func (s *StrSet) UnmarshalValue(value interface{}) (err error) {
	if s.data == nil {
		s.data = make(map[string]struct{})
	}
	var array []string
	switch value.(type) {
	case string, []byte:
		err = json.UnmarshalUseNumber(jconv.Bytes(value), &array)
	default:
		array = jconv.SliceStr(value)
	}
	for _, v := range array {
		s.data[v] = struct{}{}
	}
	return
}

// DeepCopy 返回集合的深拷贝
func (s *StrSet) DeepCopy() interface{} {
	return NewStrSetFrom(s.Slice())
}
