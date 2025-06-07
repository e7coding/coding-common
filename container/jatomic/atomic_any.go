package jatomic

import (
	"encoding/json"
	"sync/atomic"
)

// Any 是一个线程安全的任意类型值容器。
type Any struct {
	v atomic.Value
}

// NewAny 创建并返回一个可选地带初始值的 *Any。
// 如果不传参数，则 Load 后返回 nil。
func NewAny(initial ...any) *Any {
	a := &Any{}
	if len(initial) > 0 {
		a.v.Store(initial[0])
	}
	return a
}

// Load 原子地读取并返回当前值（可能为 nil）。
func (a *Any) Load() any {
	return a.v.Load()
}

// Store 原子地设置一个新值，并返回旧值。
func (a *Any) Store(new any) (old any) {
	old = a.Load()
	a.v.Store(new)
	return
}

// Clone 返回当前 Any 的浅拷贝（底层值未复制）。
func (a *Any) Clone() *Any {
	return NewAny(a.Load())
}

// MarshalJSON 实现 json.Marshaler，将底层值 JSON 序列化。
func (a *Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Load())
}

// UnmarshalJSON 实现 json.Unmarshaler，从 JSON 反序列化到底层值。
func (a *Any) UnmarshalJSON(data []byte) error {
	var x any
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	a.Store(x)
	return nil
}

// UnmarshalValue 将任意值存入底层容器。
func (a *Any) UnmarshalValue(val any) error {
	a.Store(val)
	return nil
}

// DeepCopy 返回一个深拷贝；对不可变类型（如数字/字符串）而言即为浅拷贝。
func (a *Any) DeepCopy() interface{} {
	return a.Clone()
}
