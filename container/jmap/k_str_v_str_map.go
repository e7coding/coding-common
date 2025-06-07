package jmap

import (
	"encoding/json"
	"reflect"
)

// StrStrMap 是一个简单的 string->string 映射。
type StrStrMap struct {
	data map[string]string
}

// NewStrStrMap 创建一个空的 StrStrMap。
func NewStrStrMap() *StrStrMap {
	return &StrStrMap{data: make(map[string]string)}
}

// Clone 返回当前映射的一个浅拷贝。
func (m *StrStrMap) Clone() *StrStrMap {
	cp := make(map[string]string, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &StrStrMap{data: cp}
}

// Get 获取 key 对应的值，不存在时返回空字符串。
func (m *StrStrMap) Get(key string) string {
	return m.data[key]
}

// Set 设置 key 对应的值。
func (m *StrStrMap) Set(key, val string) {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[key] = val
}

// Delete 删除指定的 key。
func (m *StrStrMap) Delete(key string) {
	delete(m.data, key)
}

// Has 判断 key 是否存在。
func (m *StrStrMap) Has(key string) bool {
	_, ok := m.data[key]
	return ok
}

// Keys 返回所有的 key 列表。
func (m *StrStrMap) Keys() []string {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有的 value 列表。
func (m *StrStrMap) Values() []string {
	vals := make([]string, 0, len(m.data))
	for _, v := range m.data {
		vals = append(vals, v)
	}
	return vals
}

// Clear 清空所有键值对。
func (m *StrStrMap) Clear() {
	m.data = make(map[string]string)
}

// Merge 将另一个映射的所有键值对合并进来，后者会覆盖前者同名键。
func (m *StrStrMap) Merge(other *StrStrMap) {
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
func (m *StrStrMap) Diff(other *StrStrMap) (added, removed, updated []string) {
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
func (m *StrStrMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

// UnmarshalJSON 支持 JSON 反序列化。
func (m *StrStrMap) UnmarshalJSON(b []byte) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return json.Unmarshal(b, &m.data)
}

// IsSubOf 判断当前映射是否是 other 的子集（键值完全一致）。
func (m *StrStrMap) IsSubOf(other *StrStrMap) bool {
	if m == other {
		return true
	}
	return reflect.DeepEqual(m, m.Clone().Intersect(other))
}

// Intersect 返回只保留同时在 m 和 other 中出现且值相等的键值对的新映射。
func (m *StrStrMap) Intersect(other *StrStrMap) *StrStrMap {
	res := NewStrStrMap()
	for k, v := range m.data {
		if ov, ok := other.data[k]; ok && ov == v {
			res.data[k] = v
		}
	}
	return res
}
