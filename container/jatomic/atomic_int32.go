package jatomic

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"sync/atomic"
)

// Int32 并发安全的 int32 类型容器。
type Int32 struct {
	v int32
}

// NewInt32 可选地以 value[0] 初始化并返回 *Int32。
func NewInt32(value ...int32) *Int32 {
	f := &Int32{}
	if len(value) > 0 {
		f.v = value[0]
	}
	return f
}

// Load 原子读取当前值。
func (f *Int32) Load() int32 {
	return atomic.LoadInt32(&f.v)
}

// Store 原子写入新值，返回旧值。
func (f *Int32) Store(val int32) (old int32) {
	return atomic.SwapInt32(&f.v, val)
}

// Add 原子地加上 delta，返回新值。
func (f *Int32) Add(delta int32) int32 {
	return atomic.AddInt32(&f.v, delta)
}

// CAS 比较并交换，返回是否成功。
func (f *Int32) CAS(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&f.v, old, new)
}

// Clone 返回当前值的拷贝容器。
func (f *Int32) Clone() *Int32 {
	return NewInt32(f.Load())
}

// String 实现 fmt.Stringer。
func (f *Int32) String() string {
	return strconv.FormatInt(int64(f.Load()), 10)
}

// MarshalJSON 实现 json.Marshaler，输出数字字面量。
func (f Int32) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(f.Load()), 10)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler，兼容 123 和 "123" 两种情况。
func (f *Int32) UnmarshalJSON(data []byte) error {
	// 尝试作为 json.Number
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		if i64, err2 := num.Int64(); err2 == nil {
			f.Store(int32(i64))
			return nil
		}
	}
	// 去除引号后解析
	s := string(bytes.Trim(data, `"`))
	i, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		f.Store(int32(i))
	}
	return err
}

// UnmarshalValue 从任意类型设置当前值，支持常见数值和字符串。
func (f *Int32) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case int:
		f.Store(int32(v))
	case int8, int16, int32, int64:
		f.Store(int32(reflect.ValueOf(v).Int()))
	case uint, uint8, uint16, uint32, uint64:
		f.Store(int32(reflect.ValueOf(v).Uint()))
	case float32, float64:
		f.Store(int32(reflect.ValueOf(v).Float()))
	case json.Number:
		if i64, err := v.Int64(); err == nil {
			f.Store(int32(i64))
		}
	case string:
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			f.Store(int32(i))
		}
	default:
		// 其他类型尝试字符串化再解析
		s := reflect.ValueOf(v).String()
		if i, err := strconv.ParseInt(s, 10, 32); err == nil {
			f.Store(int32(i))
		}
	}
	return nil
}

// DeepCopy 实现深拷贝。
func (f *Int32) DeepCopy() interface{} {
	return NewInt32(f.Load())
}
