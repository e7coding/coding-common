package jatomic

import (
	"encoding/json"
	"sync/atomic"
)

// Bool 是一个并发安全的布尔值容器。
type Bool struct {
	flag int32
}

// NewBool 返回一个可选初始值的 *Bool。
// 传入 true 时初始化为 true，否则为 false。
func NewBool(initial ...bool) *Bool {
	b := &Bool{}
	if len(initial) > 0 && initial[0] {
		b.flag = 1
	}
	return b
}

// Load 原子读取当前值。
func (b *Bool) Load() bool {
	return atomic.LoadInt32(&b.flag) == 1
}

// Store 原子写入新值，返回写入前的旧值。
func (b *Bool) Store(val bool) (old bool) {
	var x int32
	if val {
		x = 1
	}
	old = atomic.SwapInt32(&b.flag, x) == 1
	return
}

func (b *Bool) Set(value bool) (old bool) {
	if value {
		old = atomic.SwapInt32(&b.flag, 1) == 1
	} else {
		old = atomic.SwapInt32(&b.flag, 0) == 1
	}
	return
}

func (b *Bool) Val() bool {
	return atomic.LoadInt32(&b.flag) > 0
}

// CAS 执行 compare-and-swap，只有当旧值等于 old 时才会设置为 new，成功返回 true。
func (b *Bool) CAS(old, new bool) bool {
	var o, n int32
	if old {
		o = 1
	}
	if new {
		n = 1
	}
	return atomic.CompareAndSwapInt32(&b.flag, o, n)
}

// Clone 返回当前值的浅拷贝。
func (b *Bool) Clone() *Bool {
	return NewBool(b.Load())
}

// String 实现 fmt.Stringer 接口。
func (b *Bool) String() string {
	if b.Load() {
		return "true"
	}
	return "false"
}

// MarshalJSON 实现 json.Marshaler。
func (b *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Load())
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (b *Bool) UnmarshalJSON(data []byte) error {
	var v bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	b.Store(v)
	return nil
}

// UnmarshalValue 接口实现，用于接受任意类型的输入值。
func (b *Bool) UnmarshalValue(val interface{}) error {
	// 只简单地尝试转换为 bool
	v, ok := val.(bool)
	if !ok {
		return nil
	}
	b.Store(v)
	return nil
}

// DeepCopy 实现深拷贝（对布尔值而言即为浅拷贝）。
func (b *Bool) DeepCopy() interface{} {
	return b.Clone()
}
