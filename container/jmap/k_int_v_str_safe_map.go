package jmap

import (
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"sync"
)

// SafeIntStrMap 整数->字符串映射
type SafeIntStrMap struct {
	mu   sync.RWMutex
	data map[int]string
}

// NewSafeIntStrMap 创建空映射
func NewSafeIntStrMap() *SafeIntStrMap { return &SafeIntStrMap{data: make(map[int]string)} }

// NewSafeIntStrMapFrom 使用现有映射创建
func NewSafeIntStrMapFrom(d map[int]string) *SafeIntStrMap { return &SafeIntStrMap{data: d} }

// Each 遍历，f 返回 false 则中断
func (m *SafeIntStrMap) Each(f func(k int, v string) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *SafeIntStrMap) Copy() *SafeIntStrMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[int]string, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeIntStrMap{data: cp}
}

// Len 长度
func (m *SafeIntStrMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Empty 是否为空
func (m *SafeIntStrMap) Empty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Len() == 0
}

// Keys 键列表
func (m *SafeIntStrMap) Keys() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ks := make([]int, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 值列表
func (m *SafeIntStrMap) Values() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vs := make([]string, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 获取值，键不存在返回空串
func (m *SafeIntStrMap) Get(key int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// GetOK 获取值及存在标志
func (m *SafeIntStrMap) GetOK(key int) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Put 设置键值
func (m *SafeIntStrMap) Put(key int, val string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

// PutAll 批量设置
func (m *SafeIntStrMap) PutAll(d map[int]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *SafeIntStrMap) Del(key int) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *SafeIntStrMap) Pop() (int, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return 0, ""
}

// PopN 删除并返回前 N 个，N<=0 或 >Size 时返回所有
func (m *SafeIntStrMap) PopN(n int) map[int]string {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := len(m.data)
	if n <= 0 || n > t {
		n = t
	}
	res := make(map[int]string, n)
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
func (m *SafeIntStrMap) GetOrPut(key int, val string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

// Flip 交换键值
func (m *SafeIntStrMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	nm := make(map[int]string, len(m.data))
	for k, v := range m.data {
		nm[jconv.Int(v)] = jconv.String(k)
	}
	m.data = nm
}

// Prune 删除空值
func (m *SafeIntStrMap) Prune() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToStrAny 转为 map[string]interface{}
func (m *SafeIntStrMap) ToStrAny() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		res[jconv.String(k)] = v
	}
	return res
}

// UnmarshalJSON 反序列化
func (m *SafeIntStrMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var tmp map[int]string
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *SafeIntStrMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *SafeIntStrMap) DeepCopy() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[int]string, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeIntStrMap{data: cp}
}

// Sub 判断是否为子集
func (m *SafeIntStrMap) Sub(o *SafeIntStrMap) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || ov != v {
			return false
		}
	}
	return true
}

// Diff 计算新增、删除、更新键
func (m *SafeIntStrMap) Diff(o *SafeIntStrMap) (add, del, upd []int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok {
			add = append(add, k)
		} else if ov != v {
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
