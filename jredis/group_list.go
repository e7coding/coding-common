// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"github.com/e7coding/coding-common/container/jvar"
)

// 只写：Push/Insert/Set/Trim/Rem 等修改操作
type ListWriter interface {
	LPush(key string, values ...interface{}) (int64, error)
	LPushX(key string, element interface{}, elements ...interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	RPushX(key string, value interface{}) (int64, error)
	LInsert(key string, op LInsertOp, pivot, value interface{}) (int64, error)
	LSet(key string, index int64, value interface{}) (*jvar.Var, error)
	LTrim(key string, start, stop int64) error
	LRem(key string, count int64, value interface{}) (int64, error)
}

// ListReader 只读：Pop/Range/Size/Index/Blocking 等查询操作
type ListReader interface {
	LPop(key string, count ...int) (*jvar.Var, error)
	RPop(key string, count ...int) (*jvar.Var, error)
	LLen(key string) (int64, error)
	LIndex(key string, index int64) (*jvar.Var, error)
	LRange(key string, start, stop int64) (jvar.Vars, error)
	BLPop(timeout int64, keys ...string) (jvar.Vars, error)
	BRPop(timeout int64, keys ...string) (jvar.Vars, error)
	RPopLPush(source, destination string) (*jvar.Var, error)
	BRPopLPush(source, destination string, timeout int64) (*jvar.Var, error)
}

// 最终组合：
type IGroupList interface {
	ListReader
	ListWriter
}

// LInsertOp defines the operation name for function LInsert.
type LInsertOp string

const (
	LInsertBefore LInsertOp = "BEFORE"
	LInsertAfter  LInsertOp = "AFTER"
)
