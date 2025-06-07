package jmap

import (
	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// StrIntMap 字符串->整数映射
type StrIntMap struct {
	data map[string]int
}

// NewStrIntMap 创建空映射
func NewStrIntMap() *StrIntMap {
	return &StrIntMap{data: make(map[string]int)}
}

// NewFromStrIntMap 使用现有映射构造
func NewFromStrIntMap(d map[string]int) *StrIntMap {
	return &StrIntMap{data: d}
}

// Each 遍历，f 返回 false 则中断
func (m *StrIntMap) Each(f func(k string, v int) bool) {
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *StrIntMap) Copy() *StrIntMap {
	cp := make(map[string]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &StrIntMap{data: cp}
}

// Len 长度
func (m *StrIntMap) Len() int {
	return len(m.data)
}

// Empty 是否为空
func (m *StrIntMap) Empty() bool {
	return m.Len() == 0
}

// Keys 返回所有键
func (m *StrIntMap) Keys() []string {
	ks := make([]string, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 返回所有值
func (m *StrIntMap) Values() []int {
	vs := make([]int, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 返回值，键不存在返回 0
func (m *StrIntMap) Get(key string) int {
	return m.data[key]
}

// GetOK 返回值和存在标志
func (m *StrIntMap) GetOK(key string) (int, bool) {
	v, ok := m.data[key]
	return v, ok
}

// Put 设置键值
func (m *StrIntMap) Put(key string, val int) {
	m.data[key] = val
}

// PutAll 批量设置
func (m *StrIntMap) PutAll(d map[string]int) {
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *StrIntMap) Del(key string) int {
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *StrIntMap) Pop() (string, int) {
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return "", 0
}

// PopN 删除并返回前 n 对, n<=0 或 >Size 时删除所有
func (m *StrIntMap) PopN(n int) map[string]int {
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
func (m *StrIntMap) GetOrPut(key string, val int) int {
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

// Flip 交换键值，值转为字符串后作为新键
func (m *StrIntMap) Flip() {
	n := make(map[string]int, len(m.data))
	for k, v := range m.data {
		n[jconv.String(v)] = jconv.Int(k)
	}
	m.data = n
}

// Prune 删除空值
func (m *StrIntMap) Prune() {
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToVar 转为 Var
func (m *StrIntMap) ToVar(key string) *jvar.Var {
	return jvar.New(m.Get(key))
}

// ToMap 返回底层拷贝
func (m *StrIntMap) ToMap() map[string]int {
	cp := make(map[string]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return cp
}

// UnmarshalJSON 反序列化
func (m *StrIntMap) UnmarshalJSON(b []byte) error {
	var tmp map[string]int
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *StrIntMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *StrIntMap) DeepCopy() interface{} {
	cp := make(map[string]int, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &StrIntMap{data: cp}
}

// Sub 判断是否为子集
func (m *StrIntMap) Sub(o *StrIntMap) bool {
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || ov != v {
			return false
		}
	}
	return true
}

// Diff 计算新增、删除、更新键
func (m *StrIntMap) Diff(o *StrIntMap) (add, del, upd []string) {
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
