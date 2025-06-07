package jmap

import (
	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"sync"
)

// SafeStrIntMap 字符串->整数映射
type SafeStrIntMap struct {
	mu   sync.RWMutex
	data map[string]int
}

// NewSafeStrIntMap 创建空映射
func NewSafeStrIntMap() *SafeStrIntMap {
	return &SafeStrIntMap{data: make(map[string]int)}
}

// NewFromSafeStrIntMap 使用现有映射构造
func NewFromSafeStrIntMap(d map[string]int) *SafeStrIntMap {
	return &SafeStrIntMap{data: d}
}

// Each 遍历，f 返回 false 则中断
func (m *SafeStrIntMap) Each(f func(k string, v int) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *SafeStrIntMap) Copy() *SafeStrIntMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[string]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeStrIntMap{data: cp}
}

// Len 长度
func (m *SafeStrIntMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Empty 是否为空
func (m *SafeStrIntMap) Empty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Len() == 0
}

// Keys 返回所有键
func (m *SafeStrIntMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ks := make([]string, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 返回所有值
func (m *SafeStrIntMap) Values() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vs := make([]int, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 返回值，键不存在返回 0
func (m *SafeStrIntMap) Get(key string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// GetOK 返回值和存在标志
func (m *SafeStrIntMap) GetOK(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Put 设置键值
func (m *SafeStrIntMap) Put(key string, val int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

// PutAll 批量设置
func (m *SafeStrIntMap) PutAll(d map[string]int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *SafeStrIntMap) Del(key string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *SafeStrIntMap) Pop() (string, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return "", 0
}

// PopN 删除并返回前 n 对, n<=0 或 >Size 时删除所有
func (m *SafeStrIntMap) PopN(n int) map[string]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := len(m.data)
	if n <= 0 || n > t {
		n = t
	}
	res := make(map[string]int, n)
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

// GetOrPut 若存在返回, 否则设置后返回
func (m *SafeStrIntMap) GetOrPut(key string, val int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

// Flip 交换键值，值转为字符串后作为新键
func (m *SafeStrIntMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[string]int, len(m.data))
	for k, v := range m.data {
		n[jconv.String(v)] = jconv.Int(k)
	}
	m.data = n
}

// Prune 删除空值
func (m *SafeStrIntMap) Prune() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToVar 转为 Var
func (m *SafeStrIntMap) ToVar(key string) *jvar.Var {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return jvar.New(m.Get(key))
}

// ToMap 返回底层拷贝
func (m *SafeStrIntMap) ToMap() map[string]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[string]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return cp
}

// UnmarshalJSON 反序列化
func (m *SafeStrIntMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var tmp map[string]int
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *SafeStrIntMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *SafeStrIntMap) DeepCopy() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[string]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeStrIntMap{data: cp}
}

// Sub 判断是否为子集
func (m *SafeStrIntMap) Sub(o *SafeStrIntMap) bool {
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
func (m *SafeStrIntMap) Diff(o *SafeStrIntMap) (add, del, upd []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
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
