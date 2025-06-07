package jmap

import (
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"sync"
)

// SafeIntIntMap 整数映射
type SafeIntIntMap struct {
	mu   sync.RWMutex
	data map[int]int
}

// NewSafeIntIntMap 创建空映射
func NewSafeIntIntMap() *SafeIntIntMap { return &SafeIntIntMap{data: make(map[int]int)} }

// NewSafeIntIntMapFrom 使用现有映射创建
func NewSafeIntIntMapFrom(d map[int]int) *SafeIntIntMap { return &SafeIntIntMap{data: d} }

// Each 遍历，f 返回 false 则中断
func (m *SafeIntIntMap) Each(f func(k, v int) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *SafeIntIntMap) Copy() *SafeIntIntMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[int]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeIntIntMap{data: cp}
}

// Len 长度
func (m *SafeIntIntMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Empty 是否为空
func (m *SafeIntIntMap) Empty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Len() == 0
}

// Keys 键列表
func (m *SafeIntIntMap) Keys() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ks := make([]int, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 值列表
func (m *SafeIntIntMap) Values() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vs := make([]int, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 获取值，键不存在返回 0
func (m *SafeIntIntMap) Get(key int) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// GetOK 获取值和存在标志
func (m *SafeIntIntMap) GetOK(key int) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// Put 设置键值
func (m *SafeIntIntMap) Put(key, val int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

// PutAll 批量设置
func (m *SafeIntIntMap) PutAll(d map[int]int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *SafeIntIntMap) Del(key int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *SafeIntIntMap) Pop() (int, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return 0, 0
}

// PopN 删除并返回前 N 个，N<=0 或 >Size 时返回所有
func (m *SafeIntIntMap) PopN(n int) map[int]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	tot := len(m.data)
	if n <= 0 || n > tot {
		n = tot
	}
	res := make(map[int]int, n)
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
func (m *SafeIntIntMap) GetOrPut(key, val int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

// Flip 交换键值
func (m *SafeIntIntMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	nm := make(map[int]int, len(m.data))
	for k, v := range m.data {
		nm[v] = k
	}
	m.data = nm
}

// Clear 清空
func (m *SafeIntIntMap) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[int]int)
}

// Merge 合并另一个映射
func (m *SafeIntIntMap) Merge(o *SafeIntIntMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range o.data {
		m.data[k] = v
	}
}

// Prune 删除值为零或 empty 的项
func (m *SafeIntIntMap) Prune() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToStrAny 转为字符串键映射
func (m *SafeIntIntMap) ToStrAny() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		res[jconv.String(k)] = v
	}
	return res
}

// UnmarshalJSON 反序列化
func (m *SafeIntIntMap) UnmarshalJSON(b []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var tmp map[int]int
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *SafeIntIntMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *SafeIntIntMap) DeepCopy() interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(map[int]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &SafeIntIntMap{data: cp}
}

// Sub 判断是否为子集
func (m *SafeIntIntMap) Sub(o *SafeIntIntMap) bool {
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
func (m *SafeIntIntMap) Diff(o *SafeIntIntMap) (add, del, upd []int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 新增：m 有、o 没
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok {
			add = append(add, k)
		} else if ov != v {
			upd = append(upd, k)
		}
	}
	// 删除：o 有、m 没
	for k := range o.data {
		if _, ok := m.data[k]; !ok {
			del = append(del, k)
		}
	}
	return
}
