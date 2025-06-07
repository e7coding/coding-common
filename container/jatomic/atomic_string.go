package jatomic

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// String 提供原子安全的字符串操作。
type String struct {
	v atomic.Value
}

// NewString 创建并返回一个带初始值（可选）的 *String。
// 如果不传参数，则初始值为 ""，避免了后续 Load 时出现 nil。
func NewString(initial ...string) *String {
	s := &String{}
	if len(initial) > 0 {
		s.v.Store(initial[0])
	} else {
		s.v.Store("")
	}
	return s
}

// Clone 返回当前 String 的浅拷贝。
func (s *String) Clone() *String {
	return NewString(s.Load())
}

// Store 原子地设置新值，返回旧值。
func (s *String) Store(new string) (old string) {
	old = s.Load()
	s.v.Store(new)
	return
}

// Load 原子地读取当前值。
func (s *String) Load() string {
	return s.v.Load().(string)
}

// String 实现 fmt.Stringer 接口。
func (s *String) String() string {
	return s.Load()
}

// MarshalJSON 实现 json.Marshaler 接口，以 JSON 字符串格式输出。
func (s *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Load())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，从 JSON 字符串反序列化。
func (s *String) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	s.Store(str)
	return nil
}

// UnmarshalValue 将任意类型通过 fmt.Sprint 转为字符串并存储。
func (s *String) UnmarshalValue(val interface{}) error {
	s.Store(fmt.Sprint(val))
	return nil
}

// DeepCopy 返回当前对象的深拷贝（因为底层是不可变字符串，浅拷贝即深拷贝）。
func (s *String) DeepCopy() interface{} {
	return s.Clone()
}
