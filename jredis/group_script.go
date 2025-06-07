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

// IGroupScript manages redis script operations.
// Implements see jredis.GroupScript.
type IGroupScript interface {
	Eval(ctx context.Context, script string, numKeys int64, keys []string, args []interface{}) (*jvar.Var, error)
	EvalSha(ctx context.Context, sha1 string, numKeys int64, keys []string, args []interface{}) (*jvar.Var, error)
	ScriptLoad(ctx context.Context, script string) (string, error)
	ScriptExists(ctx context.Context, sha1 string, sha1s ...string) (map[string]bool, error)
	ScriptFlush(ctx context.Context, option ...ScriptFlushOption) error
	ScriptKill(ctx context.Context) error
}

// ScriptFlushOption provides options for function ScriptFlush.
type ScriptFlushOption struct {
	SYNC  bool // SYNC  flushes the cache synchronously.
	ASYNC bool // ASYNC flushes the cache asynchronously.
}
