package jset

import (
	"bytes"
	"github.com/e7coding/coding-common/jutil/jconv"
	"sync"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/text/jstr"
)

// SafeSet 是一组唯一元素的集合（非并发安全）
type SafeSet struct {
	mu   sync.RWMutex
	data map[interface{}]struct{}
}

// NewSafeSet 创建并返回一个空的集合
func NewSafeSet() *SafeSet {
	return &SafeSet{data: make(map[interface{}]struct{})}
}

// NewSafeSetFrom 从给定项（切片或单值）创建并返回一个集合
func NewSafeSetFrom(items interface{}) *SafeSet {
	s := NewSafeSet()
	for _, v := range jconv.Interfaces(items) {
		s.data[v] = struct{}{}
	}
	return s
}

// Add 向集合中添加一个或多个元素
func (s *SafeSet) Add(items ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[interface{}]struct{})
	}
	for _, v := range items {
		s.data[v] = struct{}{}
	}
}

// Contains 判断元素是否存在于集合中
func (s *SafeSet) Contains(item interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[item]
	return ok
}

// Remove 从集合中删除指定元素
func (s *SafeSet) Remove(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, item)
}

// Size 返回集合中元素的数量
func (s *SafeSet) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Clear 清空集合中所有元素
func (s *SafeSet) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[interface{}]struct{})
}

// Slice 将集合元素以切片形式返回
func (s *SafeSet) Slice() []interface{} {
	out := make([]interface{}, 0, len(s.data))
	for k := range s.data {
		out = append(out, k)
	}
	return out
}

// Iterator 遍历集合中的每个元素，f 返回 false 时停止遍历
func (s *SafeSet) Iterator(f func(interface{}) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k := range s.data {
		if !f(k) {
			return
		}
	}
}

// Union 返回多个集合的并集
func (s *SafeSet) Union(others ...*SafeSet) *SafeSet {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := NewSafeSet()
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

// Diff 返回与其他集合差集后的新集合（只包含 s 中独有的元素）
func (s *SafeSet) Diff(others ...*SafeSet) *SafeSet {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := NewSafeSet()
NEXT:
	for k := range s.data {
		for _, o := range others {
			if _, ok := o.data[k]; ok {
				continue NEXT
			}
		}
		res.data[k] = struct{}{}
	}
	return res
}

// Intersect 返回多个集合的交集
func (s *SafeSet) Intersect(others ...*SafeSet) *SafeSet {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := NewSafeSet()
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

// Complement 返回 full 与 s 的补集（full 中有而 s 中无的元素）
func (s *SafeSet) Complement(full *SafeSet) *SafeSet {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := NewSafeSet()
	for k := range full.data {
		if _, found := s.data[k]; !found {
			res.data[k] = struct{}{}
		}
	}
	return res
}

// Merge 将其他集合的元素合并到当前集合，并返回自身
func (s *SafeSet) Merge(others ...*SafeSet) *SafeSet {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, o := range others {
		for k := range o.data {
			s.data[k] = struct{}{}
		}
	}
	return s
}

// Sum 将集合元素转换为 int 并求和
func (s *SafeSet) Sum() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sum := 0
	for k := range s.data {
		sum += jconv.Int(k)
	}
	return sum
}

// Pop 随机弹出并返回一个元素；集合空时返回 nil
func (s *SafeSet) Pop() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range s.data {
		delete(s.data, k)
		return k
	}
	return nil
}

// Pops 随机弹出最多 n 个元素；n<0 或 n>Size 时返回所有元素
func (s *SafeSet) Pops(n int) []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	size := len(s.data)
	if n < 0 || n > size {
		n = size
	}
	out := make([]interface{}, 0, n)
	for k := range s.data {
		if len(out) >= n {
			break
		}
		delete(s.data, k)
		out = append(out, k)
	}
	return out
}

// Walk 对每个元素应用函数 f，然后用返回值重建集合
func (s *SafeSet) Walk(f func(interface{}) interface{}) *SafeSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	newData := make(map[interface{}]struct{}, len(s.data))
	for k := range s.data {
		newData[f(k)] = struct{}{}
	}
	s.data = newData
	return s
}

// String 返回集合的类似 JSON 数组的字符串表示
func (s *SafeSet) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var buf bytes.Buffer
	buf.WriteByte('[')
	i, n := 0, len(s.data)
	for k := range s.data {
		strs := jconv.String(k)
		if jstr.IsNumeric(strs) {
			buf.WriteString(strs)
		} else {
			buf.WriteString(`"` + jstr.QuoteMeta(strs, `"\`) + `"`)
		}
		if i < n-1 {
			buf.WriteByte(',')
		}
		i++
	}
	buf.WriteByte(']')
	return buf.String()
}

// MarshalJSON 将集合序列化为 JSON 数组
func (s *SafeSet) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.Slice())
}

// UnmarshalJSON 将 JSON 数组反序列化为集合
func (s *SafeSet) UnmarshalJSON(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[interface{}]struct{})
	}
	var arr []interface{}
	if err := json.UnmarshalUseNumber(b, &arr); err != nil {
		return err
	}
	for _, v := range arr {
		s.data[v] = struct{}{}
	}
	return nil
}
