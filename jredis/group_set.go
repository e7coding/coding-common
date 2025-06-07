// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"context"

	"github.com/e7coding/coding-common/container/jvar"
)

// SetWriter 只包含对 Redis Set 的写入/修改操作
type SetWriter interface {
	SAdd(ctx context.Context, key string, member interface{}, members ...interface{}) (int64, error)
	SRem(ctx context.Context, key string, member interface{}, members ...interface{}) (int64, error)
	SMove(ctx context.Context, source, destination string, member interface{}) (int64, error)
	SPop(ctx context.Context, key string, count ...int) (*jvar.Var, error)
	SRandMember(ctx context.Context, key string, count ...int) (*jvar.Var, error)
	SInterStore(ctx context.Context, destination string, key string, keys ...string) (int64, error)
	SUnionStore(ctx context.Context, destination, key string, keys ...string) (int64, error)
	SDiffStore(ctx context.Context, destination string, key string, keys ...string) (int64, error)
}

// SetReader 只包含对 Redis Set 的查询/只读操作
type SetReader interface {
	SIsMember(ctx context.Context, key string, member interface{}) (int64, error)
	SCard(ctx context.Context, key string) (int64, error)
	SMembers(ctx context.Context, key string) (

		jvar.Vars, error)
	SMIsMember(ctx context.Context, key string, member interface{}, members ...interface{}) ([]int, error)
	SInter(ctx context.Context, key string, keys ...string) (jvar.Vars, error)
	SUnion(ctx context.Context, key string, keys ...string) (jvar.Vars, error)
	SDiff(ctx context.Context, key string, keys ...string) (jvar.Vars, error)
}

// IGroupSet 聚合了读写接口，保持向后兼容
type IGroupSet interface {
	SetWriter
	SetReader
}
