package jatomic

import (
	"github.com/e7coding/coding-common/jutil/jconv"
	"strconv"
	"sync/atomic"
)

// Uint32 并发安全的 uint32 容器。
type Uint32 struct {
	v uint32
}

// NewUint32 可选地以 value[0] 初始化并返回 *Uint32。
func NewUint32(value ...uint32) *Uint32 {
	u := &Uint32{}
	if len(value) > 0 {
		u.Store(value[0])
	}
	return u
}

// Load 原子读取当前值。
func (u *Uint32) Load() uint32 {
	return atomic.LoadUint32(&u.v)
}

// Store 原子存储新值，返回旧值。
func (u *Uint32) Store(val uint32) (old uint32) {
	old = u.Load()
	atomic.StoreUint32(&u.v, val)
	return
}

// Swap 原子交换值，返回旧值。
func (u *Uint32) Swap(val uint32) uint32 {
	return atomic.SwapUint32(&u.v, val)
}

// Add 原子增加 delta 并返回新值。
func (u *Uint32) Add(delta uint32) uint32 {
	return atomic.AddUint32(&u.v, delta)
}
func (u *Uint32) Set(value uint32) (old uint32) {
	return atomic.SwapUint32(&u.v, value)
}

// CAS 原子比较并交换，成功时返回 true。
func (u *Uint32) CAS(old, new uint32) bool {
	return atomic.CompareAndSwapUint32(&u.v, old, new)
}

// Clone 返回当前值的浅拷贝容器。
func (u *Uint32) Clone() *Uint32 {
	return NewUint32(u.Load())
}

// String 实现 fmt.Stringer。
func (u *Uint32) String() string {
	return strconv.FormatUint(uint64(u.Load()), 10)
}

// MarshalJSON 实现 json.Marshaler。
func (u *Uint32) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u.Load()), 10)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (u *Uint32) UnmarshalJSON(data []byte) error {
	u.Set(jconv.Uint32(string(data)))
	return nil
}

// UnmarshalValue 从任意类型设置当前值。
func (u *Uint32) UnmarshalValue(val interface{}) error {
	u.Set(jconv.Uint32(val))
	return nil
}

// DeepCopy 实现深拷贝接口。
func (u *Uint32) DeepCopy() interface{} {
	if u == nil {
		return nil
	}
	return NewUint32(u.Load())
}
