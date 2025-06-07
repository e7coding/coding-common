package jatomic

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"sync/atomic"
)

// Float64 并发安全的 float64 容器。
type Float64 struct {
	v uint64
}

// NewFloat64 可选以 value[0] 初始化并返回 *Float64。
func NewFloat64(value ...float64) *Float64 {
	f := &Float64{}
	if len(value) > 0 {
		f.v = math.Float64bits(value[0])
	}
	return f
}

// Load 原子读取当前值。
func (f *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&f.v))
}

// Store 原子写入新值，返回旧值。
func (f *Float64) Store(val float64) (old float64) {
	oldBits := atomic.SwapUint64(&f.v, math.Float64bits(val))
	return math.Float64frombits(oldBits)
}

// Add 原子地加上 delta，返回新值。
func (f *Float64) Add(delta float64) (new float64) {
	for {
		oldBits := atomic.LoadUint64(&f.v)
		old := math.Float64frombits(oldBits)
		new = old + delta
		if atomic.CompareAndSwapUint64(&f.v, oldBits, math.Float64bits(new)) {
			return
		}
	}
}

// CAS 比较并交换，返回是否成功。
func (f *Float64) CAS(old, new float64) bool {
	return atomic.CompareAndSwapUint64(
		&f.v,
		math.Float64bits(old),
		math.Float64bits(new),
	)
}

// Clone 返回当前值的拷贝容器。
func (f *Float64) Clone() *Float64 {
	return NewFloat64(f.Load())
}

// String 实现 fmt.Stringer。
func (f *Float64) String() string {
	return strconv.FormatFloat(f.Load(), 'g', -1, 64)
}

// MarshalJSON 实现 json.Marshaler。
func (f *Float64) MarshalJSON() ([]byte, error) {
	// 输出数字字面量，无需引号
	return []byte(strconv.FormatFloat(f.Load(), 'g', -1, 64)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (f *Float64) UnmarshalJSON(data []byte) error {
	// 兼容带引号或不带引号的数字
	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		if fv, err2 := num.Float64(); err2 == nil {
			f.Store(fv)
			return nil
		}
	}
	// 退而尝试直接解析
	s := string(bytes.Trim(data, `"`))
	if fv, err := strconv.ParseFloat(s, 64); err == nil {
		f.Store(fv)
		return nil
	} else {
		return err
	}
}

// UnmarshalValue 从通用值设置当前值，支持多种数值类型与字符串。
func (f *Float64) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case float64:
		f.Store(v)
	case float32:
		f.Store(float64(v))
	case int:
		f.Store(float64(v))
	case int64:
		f.Store(float64(v))
	case uint:
		f.Store(float64(v))
	case uint64:
		f.Store(float64(v))
	case json.Number:
		if fv, err := v.Float64(); err == nil {
			f.Store(fv)
		}
	case string:
		if fv, err := strconv.ParseFloat(v, 64); err == nil {
			f.Store(fv)
		}
		// 其他类型忽略
	}
	return nil
}

// DeepCopy 实现深拷贝。
func (f *Float64) DeepCopy() interface{} {
	return NewFloat64(f.Load())
}
