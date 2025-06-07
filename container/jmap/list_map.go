package jmap

import (
	"bytes"
	"fmt"
	"github.com/e7coding/coding-common/container/jvar"

	"github.com/e7coding/coding-common/container/jlist"
	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
)

// ListMap 双结构：后台哈希表 + 双向链表，保持插入顺序
type ListMap struct {
	data map[interface{}]*jlist.Element
	list *jlist.List
}

type node struct {
	key, value interface{}
}

// NewListMap 创建一个空的 ListMap
func NewListMap() *ListMap {
	return &ListMap{
		data: make(map[interface{}]*jlist.Element),
		list: jlist.NewList(),
	}
}

// Clone 深拷贝 ListMap（保留顺序）
func (m *ListMap) Clone() *ListMap {
	clone := NewListMap()
	m.ForEach(func(k, v interface{}) bool {
		clone.Set(k, deepcopy.Copy(v))
		return true
	})
	return clone
}

// Clear 清空
func (m *ListMap) Clear() {
	m.data = make(map[interface{}]*jlist.Element)
	m.list = jlist.NewList()
}

// Set 添加或更新
func (m *ListMap) Set(key, val interface{}) {
	if _, exists := m.data[key]; !exists {
		m.data[key] = m.list.PushBack(&node{key, val})
	} else {
		m.data[key].Value.(*node).value = val
	}
}

// Get 返回 value，key 不存在时返回 nil
func (m *ListMap) Get(key interface{}) interface{} {
	if e, exists := m.data[key]; exists {
		return e.Value.(*node).value
	}
	return nil
}

// Search 返回 value + 是否存在
func (m *ListMap) Search(key interface{}) (interface{}, bool) {
	if e, exists := m.data[key]; exists {
		return e.Value.(*node).value, true
	}
	return nil, false
}

// Remove 删除指定 key，并返回它的 value
func (m *ListMap) Remove(key interface{}) interface{} {
	if e, exists := m.data[key]; exists {
		val := e.Value.(*node).value
		delete(m.data, key)
		m.list.Remove(e)
		return val
	}
	return nil
}

// Pop 取出并删除“任意”一对 key/value（迭代顺序中的第一个）
func (m *ListMap) Pop() (key, val interface{}) {
	for k, e := range m.data {
		key = k
		val = e.Value.(*node).value
		delete(m.data, k)
		m.list.Remove(e)
		return
	}
	return
}

// Len 返回元素个数
func (m *ListMap) Len() int {
	return len(m.data)
}

// IsEmpty 是否为空
func (m *ListMap) IsEmpty() bool {
	return m.Len() == 0
}

// Contains 判断 key 是否存在
func (m *ListMap) Contains(key interface{}) bool {
	_, ok := m.data[key]
	return ok
}

// Keys 按插入顺序返回所有 key
func (m *ListMap) Keys() []interface{} {
	out := make([]interface{}, 0, m.Len())
	m.list.ForEach(func(e *jlist.Element) bool {
		out = append(out, e.Value.(*node).key)
		return true
	})
	return out
}

// Values 按插入顺序返回所有 value
func (m *ListMap) Values() []interface{} {
	out := make([]interface{}, 0, m.Len())
	m.list.ForEach(func(e *jlist.Element) bool {
		out = append(out, e.Value.(*node).value)
		return true
	})
	return out
}

// ForEach 正序遍历
func (m *ListMap) ForEach(f func(key, value interface{}) bool) {
	m.list.ForEach(func(e *jlist.Element) bool {
		n := e.Value.(*node)
		return f(n.key, n.value)
	})
}

// ForEachReverse 逆序遍历
func (m *ListMap) ForEachReverse(f func(key, value interface{}) bool) {
	m.list.ForEachReverse(func(e *jlist.Element) bool {
		n := e.Value.(*node)
		return f(n.key, n.value)
	})
}

// FilterEmpty 删除 value 为空的条目
func (m *ListMap) FilterEmpty() {
	var toRemove []interface{}
	m.ForEach(func(k, v interface{}) bool {
		if empty.IsEmpty(v) {
			toRemove = append(toRemove, k)
		}
		return true
	})
	for _, k := range toRemove {
		m.Remove(k)
	}
}

// Merge 将 other 按顺序合并到当前，遇到已有 key 则覆盖
func (m *ListMap) Merge(other *ListMap) {
	other.ForEach(func(k, v interface{}) bool {
		m.Set(k, v)
		return true
	})
}

// Map 按插入顺序返回一个新的 map（无锁版）
func (m *ListMap) Map() map[interface{}]interface{} {
	out := make(map[interface{}]interface{}, m.Len())
	m.ForEach(func(k, v interface{}) bool {
		out[k] = v
		return true
	})
	return out
}

// MarshalJSON 自定义 JSON 序列化，按插入顺序输出
func (m *ListMap) MarshalJSON() ([]byte, error) {
	if m.data == nil {
		return []byte("null"), nil
	}
	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	first := true
	m.ForEach(func(k, v interface{}) bool {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		valBytes, err := json.Marshal(v)
		if err != nil {
			return false
		}
		buf.WriteString(fmt.Sprintf(`"%v":%s`, k, valBytes))
		return true
	})
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON 反序列化，并保持键的顺序
func (m *ListMap) UnmarshalJSON(b []byte) error {
	var tmp map[string]interface{}
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.Clear()
	for k, v := range tmp {
		m.Set(k, v)
	}
	return nil
}

// DeepCopy 返回当前 ListMap 的深拷贝
func (m *ListMap) DeepCopy() interface{} {
	clone := NewListMap()
	m.ForEach(func(k, v interface{}) bool {
		clone.Set(k, deepcopy.Copy(v))
		return true
	})
	return clone
}

// GetVar / GetVarOrSet / 等工具方法
func (m *ListMap) GetVar(key interface{}) *jvar.Var {
	return jvar.New(m.Get(key))
}
func (m *ListMap) GetOrSet(key, val interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		m.Set(key, val)
		return val
	} else {
		return v
	}
}
func (m *ListMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		v = f()
		m.Set(key, v)
	}
	return m.Get(key)
}
func (m *ListMap) SetIfNotExist(key, val interface{}) bool {
	if !m.Contains(key) {
		m.Set(key, val)
		return true
	}
	return false
}
