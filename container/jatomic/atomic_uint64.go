package jatomic

import (
	"github.com/e7coding/coding-common/jutil/jconv"
	"strconv"
	"sync/atomic"
)

// Uint64 并发安全的 uint64 容器。
type Uint64 struct {
	v uint64
}

// NewUint64 可选地以 value[0] 初始化并返回 *Uint64。
func NewUint64(value ...uint64) *Uint64 {
	u := &Uint64{}
	if len(value) > 0 {
		u.Store(value[0])
	}
	return u
}

// Load 原子读取当前值。
func (u *Uint64) Load() uint64 {
	return atomic.LoadUint64(&u.v)
}

// Store 原子存储新值，返回旧值。
func (u *Uint64) Store(val uint64) (old uint64) {
	old = u.Load()
	atomic.StoreUint64(&u.v, val)
	return
}

// Swap 原子交换值，返回旧值。
func (u *Uint64) Swap(val uint64) (old uint64) {
	return atomic.SwapUint64(&u.v, val)
}

// Add 原子增加 delta 并返回新值。
func (u *Uint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&u.v, delta)
}
func (u *Uint64) Set(value uint64) (old uint64) {
	return atomic.SwapUint64(&u.v, value)
}

// CAS 原子比较并交换，成功时返回 true。
func (u *Uint64) CAS(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&u.v, old, new)
}

// Clone 返回当前值的浅拷贝容器。
func (u *Uint64) Clone() *Uint64 {
	return NewUint64(u.Load())
}

// String 实现 fmt.Stringer。
func (u *Uint64) String() string {
	return strconv.FormatUint(u.Load(), 10)
}

// MarshalJSON 实现 json.Marshaler。
func (u *Uint64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(u.Load(), 10)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (u *Uint64) UnmarshalJSON(data []byte) error {
	s := string(data)
	// 支持带双引号的 JSON 字符串
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	u.Store(n)
	return nil
}

// UnmarshalValue 从任意类型设置当前值。
func (u *Uint64) UnmarshalValue(val interface{}) error {
	u.Set(jconv.Uint64(val))
	return nil
}

// DeepCopy 实现深拷贝接口。
func (u *Uint64) DeepCopy() interface{} {
	if u == nil {
		return nil
	}
	return NewUint64(u.Load())
}
