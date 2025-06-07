package jatomic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/e7coding/coding-common/errs/jerr"
	"sync/atomic"
)

// Bytes 是并发安全的 []byte 容器。
type Bytes struct {
	v atomic.Value // 存储 []byte
}

// NewBytes 可选地用 value[0] 创建并返回 *Bytes。
func NewBytes(value ...[]byte) *Bytes {
	b := &Bytes{}
	if len(value) > 0 && value[0] != nil {
		b.v.Store(value[0])
	}
	return b
}

// Load 原子读取当前字节切片（可能为 nil）。
func (b *Bytes) Load() []byte {
	if x := b.v.Load(); x != nil {
		return x.([]byte)
	}
	return nil
}

// Store 原子写入新字节切片 value，返回写入前的旧值。
// value 不应为 nil。
func (b *Bytes) Store(value []byte) (old []byte) {
	old = b.Load()
	b.v.Store(value)
	return
}

// Clone 返回当前字节切片的浅拷贝容器。
func (b *Bytes) Clone() *Bytes {
	old := b.Load()
	if old == nil {
		return NewBytes()
	}
	dup := make([]byte, len(old))
	copy(dup, old)
	return NewBytes(dup)
}

// String 实现 fmt.Stringer，将切片转为 string。
func (b *Bytes) String() string {
	return string(b.Load())
}

// MarshalJSON 实现 json.Marshaler，输出 Base64 编码的字符串。
func (b Bytes) MarshalJSON() ([]byte, error) {
	data := b.Load()
	enc := base64.StdEncoding.EncodeToString(data)
	return json.Marshal(enc)
}

// UnmarshalJSON 实现 json.Unmarshaler，从 Base64 解码。
// 支持带引号的 JSON 字符串。
func (b *Bytes) UnmarshalJSON(raw []byte) error {
	// 先去除引号
	s := bytes.Trim(raw, `"`)
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(dst, s)
	if err != nil {
		return jerr.WithMsgErr(err, "Base64 decode failed")
	}
	b.Store(dst[:n])
	return nil
}

// UnmarshalValue 从任意类型设置值：
// - []byte：直接存入
// - string：按原文或 Base64 解码（如果包含非可打印字符推荐 Base64）
// - 其他类型忽略
func (b *Bytes) UnmarshalValue(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		b.Store(v)
	case string:
		// 尝试按 Base64 解码
		if decoded, err := base64.StdEncoding.DecodeString(v); err == nil {
			b.Store(decoded)
		} else {
			b.Store([]byte(v))
		}
	default:
		// unsupported, ignore
	}
	return nil
}

// DeepCopy 返回深拷贝。
func (b *Bytes) DeepCopy() interface{} {
	return b.Clone()
}
