package jatomic

import (
	"encoding/json"
	"math"
	"strconv"
	"sync/atomic"
)

// Float32 并发安全的 float32 容器，
type Float32 struct {
	v uint32
}

// NewFloat32 可选地以 value[0] 初始化并返回 *Float32。
func NewFloat32(value ...float32) *Float32 {
	f := &Float32{}
	if len(value) > 0 {
		f.v = math.Float32bits(value[0])
	}
	return f
}

// Load 原子读取当前值。
func (f *Float32) Load() float32 {
	return math.Float32frombits(atomic.LoadUint32(&f.v))
}

// Store 原子写入新值，返回旧值。
func (f *Float32) Store(val float32) (old float32) {
	oldBits := atomic.SwapUint32(&f.v, math.Float32bits(val))
	return math.Float32frombits(oldBits)
}

// Add 原子地加上 delta，并返回新值。
func (f *Float32) Add(delta float32) (new float32) {
	for {
		oldBits := atomic.LoadUint32(&f.v)
		old := math.Float32frombits(oldBits)
		new = old + delta
		if atomic.CompareAndSwapUint32(
			&f.v,
			oldBits,
			math.Float32bits(new),
		) {
			return
		}
	}
}

// CAS 比较并交换，返回是否成功。
func (f *Float32) CAS(old, new float32) bool {
	return atomic.CompareAndSwapUint32(
		&f.v,
		math.Float32bits(old),
		math.Float32bits(new),
	)
}

// Clone 返回当前值的拷贝容器。
func (f *Float32) Clone() *Float32 {
	return NewFloat32(f.Load())
}

// String 实现 fmt.Stringer。
func (f *Float32) String() string {
	return strconv.FormatFloat(float64(f.Load()), 'g', -1, 32)
}

// MarshalJSON 实现 json.Marshaler。
func (f *Float32) MarshalJSON() ([]byte, error) {
	// 直接输出数字字面量
	return []byte(strconv.FormatFloat(float64(f.Load()), 'g', -1, 32)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (f *Float32) UnmarshalJSON(data []byte) error {
	// data 可能包含引号或不包含，引入 json.Unmarshal 简化处理
	var num json.Number
	if err := json.Unmarshal(data, &num); err != nil {
		// 也尝试直接去引号解析
		s := string(data)
		parsed, err2 := strconv.ParseFloat(s, 32)
		if err2 != nil {
			return err
		}
		f.Store(float32(parsed))
		return nil
	}
	parsed, err := num.Float64()
	if err != nil {
		return err
	}
	f.Store(float32(parsed))
	return nil
}

// UnmarshalValue 从通用值（string、float64、int、json.Number）设置当前值。
func (f *Float32) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case float32:
		f.Store(v)
	case float64:
		f.Store(float32(v))
	case int:
		f.Store(float32(v))
	case int64:
		f.Store(float32(v))
	case uint:
		f.Store(float32(v))
	case uint64:
		f.Store(float32(v))
	case json.Number:
		if parsed, err := v.Float64(); err == nil {
			f.Store(float32(parsed))
		}
	case string:
		if parsed, err := strconv.ParseFloat(v, 32); err == nil {
			f.Store(float32(parsed))
		}
	default:
		// unsupported types ignored
	}
	return nil
}

// DeepCopy 实现深拷贝。
func (f *Float32) DeepCopy() interface{} {
	return NewFloat32(f.Load())
}
