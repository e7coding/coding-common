// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dredis

import (
	"context"
	"github.com/coding-common/container/wmap"

	"github.com/coding-common/internal/intlog"
)

var (
	// localInstances for instance management of redis client.
	localInstances = wmap.NewSafeStrAnyMap()
)

// Instance returns an instance of redis client with specified group.
// The `name` param is unnecessary, if `name` is not passed,
// it returns a redis instance with default configuration group.
func Instance(name ...string) *Redis {
	group := DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	v := localInstances.GetOrPutFunc(group, func() interface{} {
		if config, ok := GetConfig(group); ok {
			r, err := New(config)
			if err != nil {
				intlog.Errorf(context.TODO(), `%+v`, err)
				return nil
			}
			return r
		}
		return nil
	})
	if v != nil {
		return v.(*Redis)
	}
	return nil
}
