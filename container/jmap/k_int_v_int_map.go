package jmap

import (
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// IntIntMap 整数映射
type IntIntMap struct {
	data map[int]int
}

// NewIntIntMap 创建空映射
func NewIntIntMap() *IntIntMap { return &IntIntMap{data: make(map[int]int)} }

// NewIntIntMapFrom 使用现有映射创建
func NewIntIntMapFrom(d map[int]int) *IntIntMap { return &IntIntMap{data: d} }

// Each 遍历，f 返回 false 则中断
func (m *IntIntMap) Each(f func(k, v int) bool) {
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *IntIntMap) Copy() *IntIntMap {
	cp := make(map[int]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &IntIntMap{data: cp}
}

// Len 长度
func (m *IntIntMap) Len() int { return len(m.data) }

// Empty 是否为空
func (m *IntIntMap) Empty() bool { return m.Len() == 0 }

// Keys 键列表
func (m *IntIntMap) Keys() []int {
	ks := make([]int, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 值列表
func (m *IntIntMap) Values() []int {
	vs := make([]int, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 获取值，键不存在返回 0
func (m *IntIntMap) Get(key int) int { return m.data[key] }

// GetOK 获取值和存在标志
func (m *IntIntMap) GetOK(key int) (int, bool) { v, ok := m.data[key]; return v, ok }

// Put 设置键值
func (m *IntIntMap) Put(key, val int) { m.data[key] = val }

// PutAll 批量设置
func (m *IntIntMap) PutAll(d map[int]int) {
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *IntIntMap) Del(key int) int { v := m.data[key]; delete(m.data, key); return v }

// Pop 删除并返回任意键值对
func (m *IntIntMap) Pop() (int, int) {
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return 0, 0
}

// PopN 删除并返回前 N 个，N<=0 或 >Size 时返回所有
func (m *IntIntMap) PopN(n int) map[int]int {
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
func (m *IntIntMap) GetOrPut(key, val int) int {
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

// Flip 交换键值
func (m *IntIntMap) Flip() {
	nm := make(map[int]int, len(m.data))
	for k, v := range m.data {
		nm[v] = k
	}
	m.data = nm
}

// Clear 清空
func (m *IntIntMap) Clear() { m.data = make(map[int]int) }

// Merge 合并另一个映射
func (m *IntIntMap) Merge(o *IntIntMap) {
	for k, v := range o.data {
		m.data[k] = v
	}
}

// Prune 删除值为零或 empty 的项
func (m *IntIntMap) Prune() {
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToStrAny 转为字符串键映射
func (m *IntIntMap) ToStrAny() map[string]interface{} {
	res := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		res[jconv.String(k)] = v
	}
	return res
}

// UnmarshalJSON 反序列化
func (m *IntIntMap) UnmarshalJSON(b []byte) error {
	var tmp map[int]int
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *IntIntMap) MarshalJSON() ([]byte, error) { return json.Marshal(m.data) }

// DeepCopy 深拷贝
func (m *IntIntMap) DeepCopy() interface{} {
	cp := make(map[int]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &IntIntMap{data: cp}
}

// Sub 判断是否为子集
func (m *IntIntMap) Sub(o *IntIntMap) bool {
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || ov != v {
			return false
		}
	}
	return true
}

// Diff 计算新增、删除、更新键
func (m *IntIntMap) Diff(o *IntIntMap) (add, del, upd []int) {
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
