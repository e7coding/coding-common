package jmap

import (
	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"reflect"
	"sync"
)

// SafeStrAnyMap 字符串->任意值映射
// 简化版，无并发锁
type SafeStrAnyMap struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewSafeStrAnyMap 构造空映射
func NewSafeStrAnyMap() *SafeStrAnyMap {
	return &SafeStrAnyMap{data: make(map[string]interface{})}
}

// NewSafeStrAnyMapFrom 使用现有映射创建
func NewSafeStrAnyMapFrom(d map[string]interface{}) *SafeStrAnyMap {
	return &SafeStrAnyMap{data: d}
}

// Each 遍历，f 返回 false 则中断
func (m *SafeStrAnyMap) Each(f func(k string, v interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *SafeStrAnyMap) Copy() *SafeStrAnyMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeStrAnyMap{data: cp}
}

// Len 长度
func (m *SafeStrAnyMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Empty 是否为空
func (m *SafeStrAnyMap) Empty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Len() == 0
}

// Keys 键切片
func (m *SafeStrAnyMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ks := make([]string, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 值切片
func (m *SafeStrAnyMap) Values() []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vs := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 返回值，键不存在返回 nil
func (m *SafeStrAnyMap) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// GetOK 返回值及存在标志
func (m *SafeStrAnyMap) GetOK(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Put 设置键值
func (m *SafeStrAnyMap) Put(key string, val interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

// PutAll 批量设置
func (m *SafeStrAnyMap) PutAll(d map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *SafeStrAnyMap) Del(key string) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *SafeStrAnyMap) Pop() (string, interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return "", nil
}

// PopN 删除并返回前N对，N<=0或>N时返回所有
func (m *SafeStrAnyMap) PopN(n int) map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := len(m.data)
	if n <= 0 || n > t {
		n = t
	}
	res := make(map[string]interface{}, n)
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

// GetOrPut 若存在返回，否则设置后返回
func (m *SafeStrAnyMap) GetOrPut(key string, val interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}
func (m *SafeStrAnyMap) GetOrPutFunc(key string, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithCheck(key, f)
	} else {
		return v
	}
}
func (m *SafeStrAnyMap) ByFunc(f func(m map[string]interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}
func (m *SafeStrAnyMap) Search(key string) (value interface{}, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *SafeStrAnyMap) doSetWithCheck(key string, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]interface{})
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

// Flip 交换键值，值须为可转字符串
func (m *SafeStrAnyMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		n[jconv.String(v)] = k
	}
	m.data = n
}

// Prune 删除空值
func (m *SafeStrAnyMap) Prune() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToVar 将值转为 Var
func (m *SafeStrAnyMap) ToVar(key string) *jvar.Var {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return jvar.New(m.Get(key))
}

// ToMap 转为 map[string]interface{}
func (m *SafeStrAnyMap) ToMap() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cp := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return cp
}

// UnmarshalJSON 反序列化
func (m *SafeStrAnyMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var tmp map[string]interface{}
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *SafeStrAnyMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *SafeStrAnyMap) DeepCopy() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		cp[k] = deepcopy.Copy(v)
	}
	return &SafeStrAnyMap{data: cp}
}

// Sub 判断是否为子集
func (m *SafeStrAnyMap) Sub(o *SafeStrAnyMap) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || !reflect.DeepEqual(v, ov) {
			return false
		}
	}
	return true
}

// Diff 计算新增、删除、更新键
func (m *SafeStrAnyMap) Diff(o *SafeStrAnyMap) (add, del, upd []string) {
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
