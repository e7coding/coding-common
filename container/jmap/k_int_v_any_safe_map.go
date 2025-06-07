package jmap

import (
	"github.com/e7coding/coding-common/jutil/jconv"
	"reflect"
	"sync"

	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/json"
)

// SafeIntAnyMap 整数键映射
type SafeIntAnyMap struct {
	mu   sync.RWMutex
	data map[int]interface{}
}

// NewSafeIntAnyMap 创建空映射
func NewSafeIntAnyMap() *SafeIntAnyMap { return &SafeIntAnyMap{data: make(map[int]interface{})} }

// NewSafeIntAnyMapFrom 从现有映射创建
func NewSafeIntAnyMapFrom(d map[int]interface{}) *SafeIntAnyMap { return &SafeIntAnyMap{data: d} }

// ForEach 遍历，f 返回 false 停止
func (m *SafeIntAnyMap) ForEach(f func(k int, v interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 复制映射
func (m *SafeIntAnyMap) Copy() *SafeIntAnyMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	d := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		d[k] = v
	}
	return NewSafeIntAnyMapFrom(d)
}

// Data 返回底层映射
func (m *SafeIntAnyMap) Data() map[int]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data
}

// Clone 返回浅拷贝
func (m *SafeIntAnyMap) Clone() map[int]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	d := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		d[k] = v
	}
	return d
}

// Put 设置键值
func (m *SafeIntAnyMap) Put(k int, v interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	m.data[k] = v
}

// PutAll 批量设置
func (m *SafeIntAnyMap) PutAll(d map[int]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = d
	}
	for k, v := range d {
		m.data[k] = v
	}
}

// Get 获取值
func (m *SafeIntAnyMap) Get(k int) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[k]
}

// GetOK 获取值并返回存在标志
func (m *SafeIntAnyMap) GetOK(k int) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[k]
	return v, ok
}

// Pop 随机弹出一个
func (m *SafeIntAnyMap) Pop() (int, interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return 0, nil
}

// PopN 弹出 n 个，n<0 或 >len 时弹出所有
func (m *SafeIntAnyMap) PopN(n int) map[int]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if n < 0 || n > len(m.data) {
		n = len(m.data)
	}
	res := make(map[int]interface{}, n)
	cnt := 0
	for k, v := range m.data {
		delete(m.data, k)
		res[k] = v
		cnt++
		if cnt == n {
			break
		}
	}
	return res
}

// GetOrPut 键不存在时设置并返回
func (m *SafeIntAnyMap) GetOrPut(k int, val interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.data[k]; ok {
		return existing
	}
	var v interface{}
	if fn, ok := val.(func() interface{}); ok {
		v = fn()
	} else {
		v = val
	}
	m.Put(k, v)
	return v
}

func (m *SafeIntAnyMap) ByFunc(f func(m map[int]interface{})) {
	f(m.data)
}

// Delete 删除键并返回原值
func (m *SafeIntAnyMap) Delete(k int) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.data[k]
	delete(m.data, k)
	return v
}

// Merge 合并另一个映射
func (m *SafeIntAnyMap) Merge(o *SafeIntAnyMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = o.Clone()
	}
	for k, v := range o.data {
		m.data[k] = v
	}
}

// Flip 键值互换
func (m *SafeIntAnyMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		n[jconv.Int(v)] = k
	}
	m.data = n
}

// Keys 返回键列表
func (m *SafeIntAnyMap) Keys() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ks := make([]int, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 返回值列表
func (m *SafeIntAnyMap) Values() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vs := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Has 判断是否存在
func (m *SafeIntAnyMap) Has(k int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.data[k]
	return ok
}

// Size 返回元素个数
func (m *SafeIntAnyMap) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Empty 是否为空
func (m *SafeIntAnyMap) Empty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Size() == 0
}

// Clear 清空映射
func (m *SafeIntAnyMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[int]interface{})
}

// Replace 替换数据
func (m *SafeIntAnyMap) Replace(d map[int]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = d
}

// String JSON 表示
func (m *SafeIntAnyMap) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, _ := m.MarshalJSON()
	return string(b)
}

// MarshalJSON 序列化
func (m *SafeIntAnyMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// UnmarshalJSON 反序列化
func (m *SafeIntAnyMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var d map[int]interface{}
	if err := json.UnmarshalUseNumber(b, &d); err != nil {
		return err
	}
	m.data = d
	return nil
}

// DeepCopy 深拷贝
func (m *SafeIntAnyMap) DeepCopy() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	d := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		d[k] = deepcopy.Copy(v)
	}
	return NewSafeIntAnyMapFrom(d)
}

// SubOf 是否为子集
func (m *SafeIntAnyMap) SubOf(o *SafeIntAnyMap) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m == o {
		return true
	}
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || !reflect.DeepEqual(v, ov) {
			return false
		}
	}
	return true
}

// Diff 返回新增、删除、更新键
func (m *SafeIntAnyMap) Diff(o *SafeIntAnyMap) (add, del, upd []int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok {
			add = append(add, k)
		} else if !reflect.DeepEqual(v, ov) {
			upd = append(upd, k)
		}
	}
	for k := range o.data {
		if _, ok := m.data[k]; !ok {
			del = append(del, k)
		}
	}
	return
}
