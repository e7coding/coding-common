package jmap

import (
	"github.com/e7coding/coding-common/jutil/jconv"
	"reflect"
	"sync"

	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
)

// SafeAnyMap 通用映射
type SafeAnyMap struct {
	mu   sync.RWMutex
	data map[interface{}]interface{}
}

// NewSafeAnyMap 创建空映射
func NewSafeAnyMap() *SafeAnyMap {
	return &SafeAnyMap{data: make(map[interface{}]interface{})}
}

// NewSafeAnyMapFrom 从现有 map 创建映射（不深拷贝）
func NewSafeAnyMapFrom(data map[interface{}]interface{}) *SafeAnyMap {
	return &SafeAnyMap{data: data}
}

// Iterator 遍历所有键值对，f 返回 false 时停止
func (m *SafeAnyMap) Iterator(f func(k, v interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone 返回映射副本
func (m *SafeAnyMap) Clone() *SafeAnyMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	copyMap := make(map[interface{}]interface{}, len(m.data))
	for k, v := range m.data {
		copyMap[k] = v
	}
	toMap := NewSafeAnyMapFrom(copyMap)
	return toMap
}

// Map 返回底层数据
func (m *SafeAnyMap) Map() map[interface{}]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[interface{}]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapCopy 返回数据浅拷贝
func (m *SafeAnyMap) MapCopy() map[interface{}]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	copyMap := make(map[interface{}]interface{}, len(m.data))
	for k, v := range m.data {
		copyMap[k] = v
	}
	return copyMap
}

// FilterEmpty 删除值为空的键值对
func (m *SafeAnyMap) FilterEmpty() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// FilterNil 删除值为 nil 的键值对
func (m *SafeAnyMap) FilterNil() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsNil(v) {
			delete(m.data, k)
		}
	}
}

// Set 设置键值对
func (m *SafeAnyMap) Set(key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]interface{})
	}
	m.data[key] = value
}

// Sets 批量设置键值对
func (m *SafeAnyMap) Sets(data map[interface{}]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
}

// Get 获取键对应的值
func (m *SafeAnyMap) Get(key interface{}) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val := m.data[key]
	return val

}

// Search 查找键对应的值及是否存在
func (m *SafeAnyMap) Search(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Pop 取出并删除一个随机键值对
func (m *SafeAnyMap) Pop() (key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return
}

// Pops 取出并删除指定数量键值对，size=-1 时删除所有
func (m *SafeAnyMap) Pops(size int) map[interface{}]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size < 0 || size > len(m.data) {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	res := make(map[interface{}]interface{}, size)
	count := 0
	for k, v := range m.data {
		delete(m.data, k)
		res[k] = v
		count++
		if count == size {
			break
		}
	}
	return res
}

// doSetWithCheck 键不存在时设置值并返回
func (m *SafeAnyMap) doSetWithCheck(key, value interface{}) interface{} {
	if m.data == nil {
		m.data = make(map[interface{}]interface{})
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = value
	}
	return value
}

// GetOrSet 键不存在时设置并返回值
func (m *SafeAnyMap) GetOrSet(key, value interface{}) (val interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if v, ok := m.Search(key); !ok {
		val = m.doSetWithCheck(key, value)
	} else {
		val = v
	}
	return
}

// GetOrSetFunc 键不存在时调用 f() 设置并返回
func (m *SafeAnyMap) GetOrSetFunc(key interface{}, f func() interface{}) (val interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.Search(key); !ok {
		val = m.doSetWithCheck(key, f())
	} else {
		val = v
	}
	return
}

// Remove 删除键并返回值
func (m *SafeAnyMap) Remove(key interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Merge 合并另一个映射到当前映射
func (m *SafeAnyMap) Merge(other *SafeAnyMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	for k, v := range other.data {
		m.data[k] = v
	}
}

// Flip 键值交换
func (m *SafeAnyMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[interface{}]interface{}, len(m.data))
	for k, v := range m.data {
		n[v] = k
	}
	m.data = n
}

// Keys 返回所有键
func (m *SafeAnyMap) Keys() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]interface{}, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回所有值
func (m *SafeAnyMap) Values() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vals := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		vals = append(vals, v)
	}
	return vals
}

// Contains 判断键是否存在
func (m *SafeAnyMap) Contains(key interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.data[key]
	return ok
}

// Size 返回元素个数
func (m *SafeAnyMap) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// IsEmpty 是否为空
func (m *SafeAnyMap) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Size() == 0
}

// Clear 清空映射
func (m *SafeAnyMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[interface{}]interface{})
}

// Replace 用新数据替换映射
func (m *SafeAnyMap) Replace(data map[interface{}]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = data
}

// RLockFunc 读操作回调
func (m *SafeAnyMap) RLockFunc(f func(map[interface{}]interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// LockFunc 写操作回调
func (m *SafeAnyMap) LockFunc(f func(map[interface{}]interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// String 返回 JSON 字符串
func (m *SafeAnyMap) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON 实现序列化
func (m *SafeAnyMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(jconv.Map(m.Map()))
}

// UnmarshalJSON 实现反序列化
func (m *SafeAnyMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var data map[string]interface{}
	if err := json.UnmarshalUseNumber(b, &data); err != nil {
		return err
	}
	m.data = make(map[interface{}]interface{}, len(data))
	for k, v := range data {
		m.data[k] = v
	}
	return nil
}

// DeepCopy 深拷贝
func (m *SafeAnyMap) DeepCopy() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	data := make(map[interface{}]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = deepcopy.Copy(v)
	}
	return NewSafeAnyMapFrom(data)
}

// IsSubOf 判断是否为 other 的子集
func (m *SafeAnyMap) IsSubOf(other *SafeAnyMap) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m == other {
		return true
	}
	for k, v := range m.data {
		if ov, ok := other.data[k]; !ok || !reflect.DeepEqual(v, ov) {
			return false
		}
	}
	return true
}

// Diff 返回新增、移除、更新的键
func (m *SafeAnyMap) Diff(other *SafeAnyMap) (added, removed, updated []interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if ov, ok := other.data[k]; !ok {
			added = append(added, k)
		} else if !reflect.DeepEqual(v, ov) {
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
