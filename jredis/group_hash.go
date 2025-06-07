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

// HashReader 只读：获取、枚举、查询长度、判断字段存在等
type HashReader interface {
	HGet(ctx context.Context, key, field string) (*jvar.Var, error)
	HMGet(ctx context.Context, key string, fields ...string) (jvar.Vars, error)
	HGetAll(ctx context.Context, key string) (*jvar.Var, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HVals(ctx context.Context, key string) (jvar.Vars, error)
	HExists(ctx context.Context, key, field string) (int64, error)
	HLen(ctx context.Context, key string) (int64, error)
	HStrLen(ctx context.Context, key, field string) (int64, error)
}

// HashWriter 只写：设置、删除、增量、批量写入等
type HashWriter interface {
	HSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error)
	HSetNX(ctx context.Context, key, field string, value interface{}) (int64, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HMSet(ctx context.Context, key string, fields map[string]interface{}) error
	HIncrBy(ctx context.Context, key, field string, increment int64) (int64, error)
	HIncrByFloat(ctx context.Context, key, field string, increment float64) (float64, error)
}

type IGroupHash interface {
	HashReader
	HashWriter
}
