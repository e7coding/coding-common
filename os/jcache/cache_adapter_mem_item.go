// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jcache

import (
	"github.com/e7coding/coding-common/os/jtime"
)

// IsExpired checks whether `item` is expired.
func (item *memoryDataItem) IsExpired() bool {
	// Note that it should use greater than or equal judgement here
	// imagining that the cache time is only 1 millisecond.
	return item.e < jtime.TimestampMilli()
}
