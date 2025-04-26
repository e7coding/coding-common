// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dredis

import (
	"context"

	"github.com/coding-common/container/dvar"
)

// 只写：Push/Insert/Set/Trim/Rem 等修改操作
type ListWriter interface {
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPushX(ctx context.Context, key string, element interface{}, elements ...interface{}) (int64, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPushX(ctx context.Context, key string, value interface{}) (int64, error)
	LInsert(ctx context.Context, key string, op LInsertOp, pivot, value interface{}) (int64, error)
	LSet(ctx context.Context, key string, index int64, value interface{}) (*dvar.Var, error)
	LTrim(ctx context.Context, key string, start, stop int64) error
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
}

// ListReader 只读：Pop/Range/Size/Index/Blocking 等查询操作
type ListReader interface {
	LPop(ctx context.Context, key string, count ...int) (*dvar.Var, error)
	RPop(ctx context.Context, key string, count ...int) (*dvar.Var, error)
	LLen(ctx context.Context, key string) (int64, error)
	LIndex(ctx context.Context, key string, index int64) (*dvar.Var, error)
	LRange(ctx context.Context, key string, start, stop int64) (dvar.Vars, error)
	BLPop(ctx context.Context, timeout int64, keys ...string) (dvar.Vars, error)
	BRPop(ctx context.Context, timeout int64, keys ...string) (dvar.Vars, error)
	RPopLPush(ctx context.Context, source, destination string) (*dvar.Var, error)
	BRPopLPush(ctx context.Context, source, destination string, timeout int64) (*dvar.Var, error)
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
