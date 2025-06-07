package jatomic

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"sync/atomic"
)

// Int 并发安全的整数类型容器。
type Int struct {
	v int64
}

// NewInt 可选地以 value[0] 初始化并返回 *Int。
func NewInt(value ...int) *Int {
	f := &Int{}
	if len(value) > 0 {
		f.v = int64(value[0])
	}
	return f
}

// Load 原子读取当前值。
func (f *Int) Load() int {
	return int(atomic.LoadInt64(&f.v))
}

// Store 原子写入新值，返回旧值。
func (f *Int) Store(val int) (old int) {
	oldInt := atomic.SwapInt64(&f.v, int64(val))
	return int(oldInt)
}

// Add 原子地加上 delta，返回新值。
func (f *Int) Add(delta int) int {
	return int(atomic.AddInt64(&f.v, int64(delta)))
}

// CAS 比较并交换，返回是否成功。
func (f *Int) CAS(old, new int) bool {
	return atomic.CompareAndSwapInt64(&f.v, int64(old), int64(new))
}

// Clone 返回当前值的拷贝容器。
func (f *Int) Clone() *Int {
	return NewInt(f.Load())
}

// String 实现 fmt.Stringer。
func (f *Int) String() string {
	return strconv.Itoa(f.Load())
}

// MarshalJSON 实现 json.Marshaler，输出数字字面量。
func (f Int) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(f.Load())), nil
}

// UnmarshalJSON 实现 json.Unmarshaler，兼容带引号或不带引号的数字。
func (f *Int) UnmarshalJSON(data []byte) error {
	// 先尝试标准 json.Number 解析
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		if i64, err2 := num.Int64(); err2 == nil {
			f.Store(int(i64))
			return nil
		}
	}
	// 再尝试去除引号后直接解析
	s := string(bytes.Trim(data, `"`))
	i, err := strconv.Atoi(s)
	if err == nil {
		f.Store(i)
	}
	return err
}

// UnmarshalValue 从任意类型设置当前值，支持常见数值和字符串。
func (f *Int) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case int:
		f.Store(v)
	case int8, int16, int32, int64:
		f.Store(int(reflect.ValueOf(v).Int()))
	case uint, uint8, uint16, uint32, uint64:
		f.Store(int(reflect.ValueOf(v).Uint()))
	case float32, float64:
		f.Store(int(reflect.ValueOf(v).Float()))
	case json.Number:
		if i64, err := v.Int64(); err == nil {
			f.Store(int(i64))
		}
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			f.Store(i)
		}
	}
	return nil
}

// DeepCopy 实现深拷贝。
func (f *Int) DeepCopy() interface{} {
	return NewInt(f.Load())
}
