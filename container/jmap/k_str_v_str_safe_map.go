package jmap

import (
	"encoding/json"
	"reflect"
	"sync"
)

// SafeStrStrMap 是一个简单的 string->string 映射。
type SafeStrStrMap struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewSafeStrStrMap 创建一个空的 SafeStrStrMap。
func NewSafeStrStrMap() *SafeStrStrMap {
	return &SafeStrStrMap{data: make(map[string]string)}
}

// Clone 返回当前映射的一个浅拷贝。
func (m *SafeStrStrMap) Clone() *SafeStrStrMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[string]string, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeStrStrMap{data: cp}
}

// Get 获取 key 对应的值，不存在时返回空字符串。
func (m *SafeStrStrMap) Get(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// Set 设置 key 对应的值。
func (m *SafeStrStrMap) Set(key, val string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[key] = val
}

func (m *SafeStrStrMap) SetIfNotExist(key string, value string) bool {
	if !m.Has(key) {
		m.doSetWithCheck(key, value)
		return true
	}
	return false
}

func (m *SafeStrStrMap) SetIfNotExistFunc(key string, f func() string) bool {
	if !m.Has(key) {
		m.doSetWithCheck(key, f())
		return true
	}
	return false
}
func (m *SafeStrStrMap) doSetWithCheck(key string, value string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	m.data[key] = value
	return value
}

// Delete 删除指定的 key。
func (m *SafeStrStrMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Has 判断 key 是否存在。
func (m *SafeStrStrMap) Has(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.data[key]
	return ok
}

// Keys 返回所有的 key 列表。
func (m *SafeStrStrMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有的 value 列表。
func (m *SafeStrStrMap) Values() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vals := make([]string, 0, len(m.data))
	for _, v := range m.data {
		vals = append(vals, v)
	}
	return vals
}

// Clear 清空所有键值对。
func (m *SafeStrStrMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]string)
}

// Merge 将另一个映射的所有键值对合并进来，后者会覆盖前者同名键。
func (m *SafeStrStrMap) Merge(other *SafeStrStrMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	for k, v := range other.data {
		m.data[k] = v
	}
}

// Diff 计算两个映射之间的差异：
//
//	added   - 在 m 中有但在 other 中没有的键
//	removed - 在 other 中有但在 m 中没有的键
//	updated - 在两个映射都存在但值不同的键
func (m *SafeStrStrMap) Diff(other *SafeStrStrMap) (added, removed, updated []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if ov, ok := other.data[k]; !ok {
			added = append(added, k)
		} else if ov != v {
			updated = append(updated, k)
		}
	}
	for k := range other.data {
		if _, ok := m.data[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}

// MarshalJSON 支持 JSON 序列化。
func (m *SafeStrStrMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// UnmarshalJSON 支持 JSON 反序列化。
func (m *SafeStrStrMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return json.Unmarshal(b, &m.data)
}

// IsSubOf 判断当前映射是否是 other 的子集（键值完全一致）。
func (m *SafeStrStrMap) IsSubOf(other *SafeStrStrMap) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m == other {
		return true
	}
	return reflect.DeepEqual(m, m.Clone().Intersect(other))
}

// Intersect 返回只保留同时在 m 和 other 中出现且值相等的键值对的新映射。
func (m *SafeStrStrMap) Intersect(other *SafeStrStrMap) *SafeStrStrMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := NewSafeStrStrMap()
	for k, v := range m.data {
		if ov, ok := other.data[k]; ok && ov == v {
			res.data[k] = v
		}
	}
	return res
}
