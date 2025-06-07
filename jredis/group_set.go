// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"github.com/e7coding/coding-common/container/jvar"
)

// SetWriter 只包含对 Redis Set 的写入/修改操作
type SetWriter interface {
	SAdd(key string, member interface{}, members ...interface{}) (int64, error)
	SRem(key string, member interface{}, members ...interface{}) (int64, error)
	SMove(source, destination string, member interface{}) (int64, error)
	SPop(key string, count ...int) (*jvar.Var, error)
	SRandMember(key string, count ...int) (*jvar.Var, error)
	SInterStore(destination string, key string, keys ...string) (int64, error)
	SUnionStore(destination, key string, keys ...string) (int64, error)
	SDiffStore(destination string, key string, keys ...string) (int64, error)
}

// SetReader 只包含对 Redis Set 的查询/只读操作
type SetReader interface {
	SIsMember(key string, member interface{}) (int64, error)
	SCard(key string) (int64, error)
	SMembers(key string) (jvar.Vars, error)
	SMIsMember(key string, member interface{}, members ...interface{}) ([]int, error)
	SInter(key string, keys ...string) (jvar.Vars, error)
	SUnion(key string, keys ...string) (jvar.Vars, error)
	SDiff(key string, keys ...string) (jvar.Vars, error)
}

// IGroupSet 聚合了读写接口，保持向后兼容
type IGroupSet interface {
	SetWriter
	SetReader
}
