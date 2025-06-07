// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jredis

import (
	"github.com/e7coding/coding-common/container/jvar"
)

// HashReader 只读：获取、枚举、查询长度、判断字段存在等
type HashReader interface {
	HGet(key, field string) (*jvar.Var, error)
	HMGet(key string, fields ...string) (jvar.Vars, error)
	HGetAll(key string) (*jvar.Var, error)
	HKeys(key string) ([]string, error)
	HVals(key string) (jvar.Vars, error)
	HExists(key, field string) (int64, error)
	HLen(key string) (int64, error)
	HStrLen(key, field string) (int64, error)
}

// HashWriter 只写：设置、删除、增量、批量写入等
type HashWriter interface {
	HSet(key string, fields map[string]interface{}) (int64, error)
	HSetNX(key, field string, value interface{}) (int64, error)
	HDel(key string, fields ...string) (int64, error)
	HMSet(key string, fields map[string]interface{}) error
	HIncrBy(key, field string, increment int64) (int64, error)
	HIncrByFloat(key, field string, increment float64) (float64, error)
}

type IGroupHash interface {
	HashReader
	HashWriter
}
