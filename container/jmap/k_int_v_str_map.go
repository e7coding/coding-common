package jmap

import (
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
)

// IntStrMap 整数->字符串映射
type IntStrMap struct {
	data map[int]string
}

// NewIntStrMap 创建空映射
func NewIntStrMap() *IntStrMap { return &IntStrMap{data: make(map[int]string)} }

// NewIntStrMapFrom 使用现有映射创建
func NewIntStrMapFrom(d map[int]string) *IntStrMap { return &IntStrMap{data: d} }

// Each 遍历，f 返回 false 则中断
func (m *IntStrMap) Each(f func(k int, v string) bool) {
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *IntStrMap) Copy() *IntStrMap {
	cp := make(map[int]string, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &IntStrMap{data: cp}
}

// Len 长度
func (m *IntStrMap) Len() int { return len(m.data) }

// Empty 是否为空
func (m *IntStrMap) Empty() bool { return m.Len() == 0 }

// Keys 键列表
func (m *IntStrMap) Keys() []int {
	ks := make([]int, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 值列表
func (m *IntStrMap) Values() []string {
	vs := make([]string, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 获取值，键不存在返回空串
func (m *IntStrMap) Get(key int) string { return m.data[key] }

// GetOK 获取值及存在标志
func (m *IntStrMap) GetOK(key int) (string, bool) { v, ok := m.data[key]; return v, ok }

// Put 设置键值
func (m *IntStrMap) Put(key int, val string) { m.data[key] = val }

// PutAll 批量设置
func (m *IntStrMap) PutAll(d map[int]string) {
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *IntStrMap) Del(key int) string {
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *IntStrMap) Pop() (int, string) {
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return 0, ""
}

// PopN 删除并返回前 N 个，N<=0 或 >Size 时返回所有
func (m *IntStrMap) PopN(n int) map[int]string {
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
func (m *IntStrMap) GetOrPut(key int, val string) string {
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

// Flip 交换键值
func (m *IntStrMap) Flip() {
	nm := make(map[int]string, len(m.data))
	for k, v := range m.data {
		nm[jconv.Int(v)] = jconv.String(k)
	}
	m.data = nm
}

// Prune 删除空值
func (m *IntStrMap) Prune() {
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToStrAny 转为 map[string]interface{}
func (m *IntStrMap) ToStrAny() map[string]interface{} {
	res := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		res[jconv.String(k)] = v
	}
	return res
}

// UnmarshalJSON 反序列化
func (m *IntStrMap) UnmarshalJSON(b []byte) error {
	var tmp map[int]string
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *IntStrMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *IntStrMap) DeepCopy() interface{} {
	cp := make(map[int]string, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &IntStrMap{data: cp}
}

// Sub 判断是否为子集
func (m *IntStrMap) Sub(o *IntStrMap) bool {
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || ov != v {
			return false
		}
	}
	return true
}

// Diff 计算新增、删除、更新键
func (m *IntStrMap) Diff(o *IntStrMap) (add, del, upd []int) {
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
