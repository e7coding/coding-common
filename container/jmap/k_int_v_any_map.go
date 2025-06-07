package jmap

import (
	"github.com/e7coding/coding-common/jutil/jconv"
	"reflect"

	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/json"
)

// IntAnyMap 整数键映射
type IntAnyMap struct{ data map[int]interface{} }

// NewIntAnyMap 创建空映射
func NewIntAnyMap() *IntAnyMap { return &IntAnyMap{data: make(map[int]interface{})} }

// NewIntAnyMapFrom 从现有映射创建
func NewIntAnyMapFrom(d map[int]interface{}) *IntAnyMap { return &IntAnyMap{data: d} }

// ForEach 遍历，f 返回 false 停止
func (m *IntAnyMap) ForEach(f func(k int, v interface{}) bool) {
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 复制映射
func (m *IntAnyMap) Copy() *IntAnyMap {
	d := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		d[k] = v
	}
	return NewIntAnyMapFrom(d)
}

// Data 返回底层映射
func (m *IntAnyMap) Data() map[int]interface{} { return m.data }

// Clone 返回浅拷贝
func (m *IntAnyMap) Clone() map[int]interface{} {
	d := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		d[k] = v
	}
	return d
}

// Put 设置键值
func (m *IntAnyMap) Put(k int, v interface{}) {
	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	m.data[k] = v
}

// PutAll 批量设置
func (m *IntAnyMap) PutAll(d map[int]interface{}) {
	if m.data == nil {
		m.data = d
	}
	for k, v := range d {
		m.data[k] = v
	}
}

// Get 获取值
func (m *IntAnyMap) Get(k int) interface{} { return m.data[k] }

// GetOK 获取值并返回存在标志
func (m *IntAnyMap) GetOK(k int) (interface{}, bool) {
	v, ok := m.data[k]
	return v, ok
}

// Pop 随机弹出一个
func (m *IntAnyMap) Pop() (int, interface{}) {
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return 0, nil
}

// PopN 弹出 n 个，n<0 或 >len 时弹出所有
func (m *IntAnyMap) PopN(n int) map[int]interface{} {
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
func (m *IntAnyMap) GetOrPut(k int, val interface{}) interface{} {
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

// Delete 删除键并返回原值
func (m *IntAnyMap) Delete(k int) interface{} {
	v := m.data[k]
	delete(m.data, k)
	return v
}

// Merge 合并另一个映射
func (m *IntAnyMap) Merge(o *IntAnyMap) {
	if m.data == nil {
		m.data = o.Clone()
	}
	for k, v := range o.data {
		m.data[k] = v
	}
}

// Flip 键值互换
func (m *IntAnyMap) Flip() {
	n := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		n[jconv.Int(v)] = k
	}
	m.data = n
}

// Keys 返回键列表
func (m *IntAnyMap) Keys() []int {
	ks := make([]int, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 返回值列表
func (m *IntAnyMap) Values() []interface{} {
	vs := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Has 判断是否存在
func (m *IntAnyMap) Has(k int) bool { _, ok := m.data[k]; return ok }

// Len 返回元素个数
func (m *IntAnyMap) Len() int { return len(m.data) }

// Empty 是否为空
func (m *IntAnyMap) Empty() bool { return m.Len() == 0 }

// Clear 清空映射
func (m *IntAnyMap) Clear() { m.data = make(map[int]interface{}) }

// Replace 替换数据
func (m *IntAnyMap) Replace(d map[int]interface{}) { m.data = d }

// String JSON 表示
func (m *IntAnyMap) String() string { b, _ := m.MarshalJSON(); return string(b) }

// MarshalJSON 序列化
func (m *IntAnyMap) MarshalJSON() ([]byte, error) { return json.Marshal(m.data) }

// UnmarshalJSON 反序列化
func (m *IntAnyMap) UnmarshalJSON(b []byte) error {
	var d map[int]interface{}
	if err := json.UnmarshalUseNumber(b, &d); err != nil {
		return err
	}
	m.data = d
	return nil
}

// DeepCopy 深拷贝
func (m *IntAnyMap) DeepCopy() interface{} {
	d := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		d[k] = deepcopy.Copy(v)
	}
	return NewIntAnyMapFrom(d)
}

// SubOf 是否为子集
func (m *IntAnyMap) SubOf(o *IntAnyMap) bool {
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
func (m *IntAnyMap) Diff(o *IntAnyMap) (add, del, upd []int) {
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
