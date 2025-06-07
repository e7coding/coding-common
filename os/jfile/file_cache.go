// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jfile

import (
	"github.com/e7coding/coding-common/errs/jerr"
	"time"

	"github.com/e7coding/coding-common/internal/intlog"
	"github.com/e7coding/coding-common/os/jcache"
	"github.com/e7coding/coding-common/os/jfsnotify"
)

const (
	defaultCacheDuration  = "1m" // defaultCacheExpire is the expire time for file content caching in seconds.
	commandEnvKeyForCache = "coding"
)

var (
	// Default expire time for file content caching.
	cacheDuration = getCacheDuration()

	// internalCache is the memory cache for internal usage.
	internalCache = jcache.New()
)

func getCacheDuration() time.Duration {
	d, err := time.ParseDuration(defaultCacheDuration)
	if err != nil {
		panic(jerr.WithMsgErrF(
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
		expire   = cacheDuration
		cacheKey = commandEnvKeyForCache + path
	)

	if len(duration) > 0 {
		expire = duration[0]
	}
	r, _ := internalCache.GetOrSetFuncLock(cacheKey, func() (interface{}, error) {
		b := GetBytes(path)
		if b != nil {
			// Adding this `path` to gfsnotify,
			// it will clear its cache if there's any changes of the file.
			_, _ = jfsnotify.Add(path, func(event *jfsnotify.Event) {
				_, err := internalCache.Remove(cacheKey)
				if err != nil {
					intlog.Errorf(`%+v`, err)
				}
				jfsnotify.Exit()
			})
		}
		return b, nil
	}, expire)
	if r != nil {
		return r.Bytes()
	}
	return nil
}
