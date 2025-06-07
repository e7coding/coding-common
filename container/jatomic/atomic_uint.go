package jatomic

import (
	"github.com/e7coding/coding-common/jutil/jconv"
	"strconv"
	"sync/atomic"
)

// Uint 是一个并发安全的 uint 容器。
type Uint struct {
	v uint64
}

// NewUint 可选地以 value[0] 初始化并返回 *Uint。
func NewUint(value ...uint) *Uint {
	u := &Uint{}
	if len(value) > 0 {
		u.Store(value[0])
	}
	return u
}

// Load 原子读取当前值。
func (u *Uint) Load() uint {
	return uint(atomic.LoadUint64(&u.v))
}

// Store 原子存储新值，返回旧值。
func (u *Uint) Store(val uint) (old uint) {
	old = u.Load()
	atomic.StoreUint64(&u.v, uint64(val))
	return
}

// Swap 原子交换值，返回旧值。
func (u *Uint) Swap(val uint) (old uint) {
	old = u.Load()
	atomic.SwapUint64(&u.v, uint64(val))
	return
}
func (u *Uint) Set(value uint) (old uint) {
	return uint(atomic.SwapUint64(&u.v, uint64(value)))
}

// Add 原子增加 delta 并返回新值。
func (u *Uint) Add(delta uint) uint {
	return uint(atomic.AddUint64(&u.v, uint64(delta)))
}

// CAS 原子比较并交换，成功时返回 true。
func (u *Uint) CAS(old, new uint) bool {
	return atomic.CompareAndSwapUint64(&u.v, uint64(old), uint64(new))
}

// Clone 返回当前值的浅拷贝容器。
func (u *Uint) Clone() *Uint {
	return NewUint(u.Load())
}

// String 实现 fmt.Stringer。
func (u *Uint) String() string {
	return strconv.FormatUint(uint64(u.Load()), 10)
}

// MarshalJSON 实现 json.Marshaler。
func (u *Uint) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u.Load()), 10)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (u *Uint) UnmarshalJSON(data []byte) error {
	u.Set(jconv.Uint(string(data)))
	return nil
}

// UnmarshalValue 从任意类型设置当前值。
func (u *Uint) UnmarshalValue(val interface{}) error {
	u.Set(jconv.Uint(val))
	return nil
}

// DeepCopy 实现深拷贝接口。
func (u *Uint) DeepCopy() interface{} {
	if u == nil {
		return nil
	}
	return NewUint(u.Load())
}
