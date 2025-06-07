package jatomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

// Byte 是一个并发安全的 uint8（byte）容器。
type Byte struct {
	v int32
}

// NewByte 可选地用初始值 value[0] 创建并返回 *Byte。
func NewByte(value ...byte) *Byte {
	b := &Byte{}
	if len(value) > 0 {
		b.v = int32(value[0])
	}
	return b
}

// Load 原子读取当前值。
func (b *Byte) Load() byte {
	return byte(atomic.LoadInt32(&b.v))
}

// Store 原子写入新值，返回写入前的旧值。
func (b *Byte) Store(x byte) (old byte) {
	old = byte(atomic.SwapInt32(&b.v, int32(x)))
	return
}

// Add 原子将 delta 加到当前值上，返回新的值。（溢出到 0-255 周期内）
func (b *Byte) Add(delta byte) byte {
	return byte(atomic.AddInt32(&b.v, int32(delta)))
}

// CAS 原子执行比较并交换，只有当旧值等于 old 时才会设置为 new，成功返回 true。
func (b *Byte) CAS(old, new byte) bool {
	return atomic.CompareAndSwapInt32(&b.v, int32(old), int32(new))
}

// Clone 返回当前值的浅拷贝。
func (b *Byte) Clone() *Byte {
	return NewByte(b.Load())
}

// String 实现 fmt.Stringer 接口，以十进制形式返回。
func (b *Byte) String() string {
	return strconv.FormatUint(uint64(b.Load()), 10)
}

// MarshalJSON 实现 json.Marshaler，将值按数字输出。
func (b Byte) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Load())
}

// UnmarshalJSON 实现 json.Unmarshaler，从数字或带引号的数字反序列化。
func (b *Byte) UnmarshalJSON(data []byte) error {
	// 既支持 "123" 也支持 123
	s := string(data)
	// Trim 引号
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	u, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	b.Store(byte(u))
	return nil
}

// UnmarshalValue 从任意类型值解析并存入，未能解析则忽略。
func (b *Byte) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case byte:
		b.Store(v)
	case int:
		b.Store(byte(v))
	case float64:
		b.Store(byte(v))
	case string:
		u, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return err
		}
		b.Store(byte(u))
	default:
		// 其他类型忽略
	}
	return nil
}

// DeepCopy 实现深拷贝。
func (b *Byte) DeepCopy() interface{} {
	return b.Clone()
}
