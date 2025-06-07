package jmap

import (
	"github.com/e7coding/coding-common/container/jvar"
	"github.com/e7coding/coding-common/internal/deepcopy"
	"github.com/e7coding/coding-common/internal/empty"
	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"reflect"
)

// StrAnyMap 字符串->任意值映射
// 简化版，无并发锁
type StrAnyMap struct {
	data map[string]interface{}
}

// NewStrAnyMap 构造空映射
func NewStrAnyMap() *StrAnyMap {
	return &StrAnyMap{data: make(map[string]interface{})}
}

// NewStrAnyMapFrom 使用现有映射创建
func NewStrAnyMapFrom(d map[string]interface{}) *StrAnyMap {
	return &StrAnyMap{data: d}
}

// Each 遍历，f 返回 false 则中断
func (m *StrAnyMap) Each(f func(k string, v interface{}) bool) {
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Copy 返回浅拷贝
func (m *StrAnyMap) Copy() *StrAnyMap {
	cp := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return &StrAnyMap{data: cp}
}

// Len 长度
func (m *StrAnyMap) Len() int { return len(m.data) }

// Empty 是否为空
func (m *StrAnyMap) Empty() bool { return m.Len() == 0 }

// Keys 键切片
func (m *StrAnyMap) Keys() []string {
	ks := make([]string, 0, len(m.data))
	for k := range m.data {
		ks = append(ks, k)
	}
	return ks
}

// Values 值切片
func (m *StrAnyMap) Values() []interface{} {
	vs := make([]interface{}, 0, len(m.data))
	for _, v := range m.data {
		vs = append(vs, v)
	}
	return vs
}

// Get 返回值，键不存在返回 nil
func (m *StrAnyMap) Get(key string) interface{} {
	return m.data[key]
}

// GetOK 返回值及存在标志
func (m *StrAnyMap) GetOK(key string) (interface{}, bool) {
	v, ok := m.data[key]
	return v, ok
}

// Put 设置键值
func (m *StrAnyMap) Put(key string, val interface{}) {
	m.data[key] = val
}

// PutAll 批量设置
func (m *StrAnyMap) PutAll(d map[string]interface{}) {
	for k, v := range d {
		m.data[k] = v
	}
}

// Del 删除并返回原值
func (m *StrAnyMap) Del(key string) interface{} {
	v := m.data[key]
	delete(m.data, key)
	return v
}

// Pop 删除并返回任意键值对
func (m *StrAnyMap) Pop() (string, interface{}) {
	for k, v := range m.data {
		delete(m.data, k)
		return k, v
	}
	return "", nil
}

// PopN 删除并返回前N对，N<=0或>N时返回所有
func (m *StrAnyMap) PopN(n int) map[string]interface{} {
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
func (m *StrAnyMap) GetOrPut(key string, val interface{}) interface{} {
	if v, ok := m.GetOK(key); ok {
		return v
	}
	m.Put(key, val)
	return val
}

func (m *StrAnyMap) GetOrPutFunc(key string, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithCheck(key, f)
	} else {
		return v
	}
}
func (m *StrAnyMap) Search(key string) (value interface{}, found bool) {
	if m.data != nil {
		value, found = m.data[key]
	}
	return
}

func (m *StrAnyMap) doSetWithCheck(key string, value interface{}) interface{} {
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
func (m *StrAnyMap) Flip() {
	n := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		n[jconv.String(v)] = k
	}
	m.data = n
}

// Prune 删除空值
func (m *StrAnyMap) Prune() {
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
}

// ToVar 将值转为 Var
func (m *StrAnyMap) ToVar(key string) *jvar.Var {
	return jvar.New(m.Get(key))
}

// ToMap 转为 map[string]interface{}
func (m *StrAnyMap) ToMap() map[string]interface{} {
	cp := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		cp[k] = v
	}
	return cp
}

// UnmarshalJSON 反序列化
func (m *StrAnyMap) UnmarshalJSON(b []byte) error {
	var tmp map[string]interface{}
	if err := json.UnmarshalUseNumber(b, &tmp); err != nil {
		return err
	}
	m.data = tmp
	return nil
}

// MarshalJSON 序列化
func (m *StrAnyMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

// DeepCopy 深拷贝
func (m *StrAnyMap) DeepCopy() interface{} {
	cp := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		cp[k] = deepcopy.Copy(v)
	}
	return &StrAnyMap{data: cp}
}

// Sub 判断是否为子集
func (m *StrAnyMap) Sub(o *StrAnyMap) bool {
	for k, v := range m.data {
		if ov, ok := o.data[k]; !ok || !reflect.DeepEqual(v, ov) {
			return false
		}
	}
	return true
}

// Diff 计算新增、删除、更新键
func (m *StrAnyMap) Diff(o *StrAnyMap) (add, del, upd []string) {
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
