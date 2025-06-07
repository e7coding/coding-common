package jatomic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/e7coding/coding-common/internal/deepcopy"
)

// Interface 是一个并发安全的任意值容器（原 Interface）。
type Interface struct {
	v atomic.Value
}

// NewInterface 可选地以 value[0] 初始化并返回 *Interface。
func NewInterface(value ...interface{}) *Interface {
	a := &Interface{}
	if len(value) > 0 && value[0] != nil {
		a.v.Store(value[0])
	}
	return a
}

// Load 原子读取当前存储的值。
func (a *Interface) Load() interface{} {
	return a.v.Load()
}

// Store 原子写入新值，返回旧值。
func (a *Interface) Store(val interface{}) (old interface{}) {
	old = a.Load()
	a.v.Store(val)
	return
}

// Clone 返回当前值的浅拷贝容器。
func (a *Interface) Clone() *Interface {
	return NewInterface(a.Load())
}

// String 实现 fmt.Stringer。
// 如果底层就是字符串，直接返回；否则尝试 JSON 序列化。
func (a *Interface) String() string {
	if s, ok := a.Load().(string); ok {
		return s
	}
	b, err := json.Marshal(a.Load())
	if err != nil {
		return fmt.Sprint(a.Load())
	}
	return string(b)
}

// MarshalJSON 实现 json.Marshaler，直接序列化底层值。
func (a *Interface) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Load())
}

// UnmarshalJSON 实现 json.Unmarshaler，使用 UseNumber 保持数值精度。
func (a *Interface) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var x interface{}
	if err := dec.Decode(&x); err != nil {
		return err
	}
	a.Store(x)
	return nil
}

// UnmarshalValue 从任意类型设置当前值。
func (a *Interface) UnmarshalValue(val interface{}) error {
	a.Store(val)
	return nil
}

// DeepCopy 深拷贝底层值（利用 coding-common 内部 deepcopy）。
func (a *Interface) DeepCopy() interface{} {
	if a == nil {
		return nil
	}
	return NewInterface(deepcopy.Copy(a.Load()))
}
