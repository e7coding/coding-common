// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jres

import "github.com/e7coding/coding-common/container/jmap"

const (
	// DefaultName default group name for instance usage.
	DefaultName = "default"
)

var (
	// Instances map.
	instances = jmap.NewSafeStrAnyMap()
)

// Instance returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Instance(name ...string) *Resource {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrPutFunc(key, func() interface{} {
		return New()
	}).(*Resource)
}
