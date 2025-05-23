// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfile

import (
	"context"
	"github.com/coding-common/errs/werr"
	"time"

	"github.com/coding-common/internal/intlog"
	"github.com/coding-common/os/dcache"
	"github.com/coding-common/os/dfsnotify"
)

const (
	defaultCacheDuration  = "1m" // defaultCacheExpire is the expire time for file content caching in seconds.
	commandEnvKeyForCache = "coding"
)

var (
	// Default expire time for file content caching.
	cacheDuration = getCacheDuration()

	// internalCache is the memory cache for internal usage.
	internalCache = dcache.New()
)

func getCacheDuration() time.Duration {
	d, err := time.ParseDuration(defaultCacheDuration)
	if err != nil {
		panic(werr.WithMsgErrF(
			err,
			`error  time duration`,
		))
	}
	return d
}

// GetContentsWithCache returns string content of given file by `path` from cache.
// If there's no content in the cache, it will read it from disk file specified by `path`.
// The parameter `expire` specifies the caching time for this file content in seconds.
func GetContentsWithCache(path string, duration ...time.Duration) string {
	return string(GetBytesWithCache(path, duration...))
}

// GetBytesWithCache returns []byte content of given file by `path` from cache.
// If there's no content in the cache, it will read it from disk file specified by `path`.
// The parameter `expire` specifies the caching time for this file content in seconds.
func GetBytesWithCache(path string, duration ...time.Duration) []byte {
	var (
		ctx      = context.Background()
		expire   = cacheDuration
		cacheKey = commandEnvKeyForCache + path
	)

	if len(duration) > 0 {
		expire = duration[0]
	}
	r, _ := internalCache.GetOrSetFuncLock(ctx, cacheKey, func(ctx context.Context) (interface{}, error) {
		b := GetBytes(path)
		if b != nil {
			// Adding this `path` to gfsnotify,
			// it will clear its cache if there's any changes of the file.
			_, _ = dfsnotify.Add(path, func(event *dfsnotify.Event) {
				_, err := internalCache.Remove(ctx, cacheKey)
				if err != nil {
					intlog.Errorf(ctx, `%+v`, err)
				}
				dfsnotify.Exit()
			})
		}
		return b, nil
	}, expire)
	if r != nil {
		return r.Bytes()
	}
	return nil
}
