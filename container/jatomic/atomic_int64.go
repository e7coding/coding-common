package jatomic

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"sync/atomic"
)

// Int64 并发安全的 int64 类型容器。
type Int64 struct {
	v int64
}

// NewInt64 可选地以 value[0] 初始化并返回 *Int64。
func NewInt64(value ...int64) *Int64 {
	f := &Int64{}
	if len(value) > 0 {
		f.v = value[0]
	}
	return f
}

// Load 原子读取当前值。
func (f *Int64) Load() int64 {
	return atomic.LoadInt64(&f.v)
}

// Store 原子写入新值，返回旧值。
func (f *Int64) Store(val int64) (old int64) {
	return atomic.SwapInt64(&f.v, val)
}

// Add 原子地加上 delta，返回新值。
func (f *Int64) Add(delta int64) int64 {
	return atomic.AddInt64(&f.v, delta)
}

// CAS 比较并交换，返回是否成功。
func (f *Int64) CAS(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&f.v, old, new)
}

// Clone 返回当前值的拷贝容器。
func (f *Int64) Clone() *Int64 {
	return NewInt64(f.Load())
}

// String 实现 fmt.Stringer。
func (f *Int64) String() string {
	return strconv.FormatInt(f.Load(), 10)
}

// MarshalJSON 实现 json.Marshaler，输出数字字面量。
func (f *Int64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(f.Load(), 10)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler，兼容 123 和 "123" 两种情况。
func (f *Int64) UnmarshalJSON(data []byte) error {
	// 尝试作为 json.Number
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		if i64, err2 := num.Int64(); err2 == nil {
			f.Store(i64)
			return nil
		}
	}
	// 去除引号后解析
	s := string(bytes.Trim(data, `"`))
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		f.Store(i)
	}
	return err
}

// UnmarshalValue 从任意类型设置当前值，支持常见数值和字符串。
func (f *Int64) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case int:
		f.Store(int64(v))
	case int8, int16, int32, int64:
		f.Store(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		f.Store(int64(reflect.ValueOf(v).Uint()))
	case float32, float64:
		f.Store(int64(reflect.ValueOf(v).Float()))
	case json.Number:
		if i64, err := v.Int64(); err == nil {
			f.Store(i64)
		}
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.Store(i)
		}
	default:
		s := reflect.ValueOf(v).String()
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			f.Store(i)
		}
	}
	return nil
}

// DeepCopy 实现深拷贝。
func (f *Int64) DeepCopy() interface{} {
	return NewInt64(f.Load())
}
