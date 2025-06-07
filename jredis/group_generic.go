// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"time"
)

type IGroupGeneric interface {
	KeyOps
	ScanOps
	FlushOps
	ExpireOps
	CopyOps
	Type(key string) (string, error)
	Unlink(keys ...string) (int64, error)
	DBSize() (int64, error)
}

// KeyOps 只关心 key 的基本操作
type KeyOps interface {
	Exists(keys ...string) (int64, error)
	Del(keys ...string) (int64, error)
	Rename(key, newKey string) error
	RenameNX(key, newKey string) (int64, error)
	Move(key string, db int) (int64, error)
	RandomKey() (string, error)
}

// ScanOps 只关心遍历和扫描
type ScanOps interface {
	Keys(pattern string) ([]string, error)
	Scan(cursor uint64, option ...ScanOption) (uint64, []string, error)
}

// FlushOps 只关心清库操作
type FlushOps interface {
	FlushDB(option ...FlushOp) error
	FlushAll(option ...FlushOp) error
}

// TTL 和过期相关
type ExpireOps interface {
	Expire(key string, seconds int64, option ...ExpireOption) (int64, error)
	ExpireAt(key string, when time.Time, option ...ExpireOption) (int64, error)
	TTL(key string) (int64, error)
	Persist(key string) (int64, error)
	PExpire(key string, ms int64, option ...ExpireOption) (int64, error)
	PExpireAt(key string, when time.Time, option ...ExpireOption) (int64, error)
	PTTL(key string) (int64, error)
}

// CopyOps 只关心复制命令
type CopyOps interface {
	Copy(source, dest string, option ...CopyOption) (int64, error)
}

// CopyOption provides options for function Copy.
type CopyOption struct {
	DB      int  // DB option allows specifying an alternative logical database index for the destination key.
	REPLACE bool // REPLACE option removes the destination key before copying the value to it.
}

type FlushOp string

const (
	FlushAsync FlushOp = "ASYNC" // ASYNC: flushes the databases asynchronously
	FlushSync  FlushOp = "SYNC"  // SYNC: flushes the databases synchronously
)

// ExpireOption provides options for function Expire.
type ExpireOption struct {
	NX bool // NX -- Set expiry only when the key has no expiry
	XX bool // XX -- Set expiry only when the key has an existing expiry
	GT bool // GT -- Set expiry only when the new expiry is greater than current one
	LT bool // LT -- Set expiry only when the new expiry is less than current one
}

// ScanOption provides options for function Scan.
type ScanOption struct {
	Match string // Match -- Specifies a glob-style pattern for filtering keys.
	Count int    // Count -- Suggests the number of keys to return per scan.
	Type  string // Type -- Filters keys by their data type. Valid types are "string", "list", "set", "zset", "hash", and "stream".
}

// ToUsedOption converts fields in ScanOption with zero values to nil. Only fields with values are retained.
func (so *ScanOption) ToUsedOption() *ScanOption {
	usedOption := &ScanOption{}
	if so.Match != "" {
		usedOption.Match = so.Match
	}
	if so.Count != 0 {
		usedOption.Count = so.Count
	}
	if so.Type != "" {
		usedOption.Type = so.Type
	}

	return usedOption
}
